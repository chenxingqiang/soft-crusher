package deployment

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/chenxingqiang/soft-crusher/internal/config"
	"github.com/chenxingqiang/soft-crusher/pkg/logging"
	"go.uber.org/zap"
)

type AWSDeployer struct {
	cfg       *config.Config
	ec2Client *ec2.EC2
	ecsClient *ecs.ECS
}

func NewAWSDeployer(cfg *config.Config) (*AWSDeployer, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(cfg.GetCloudRegion()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	return &AWSDeployer{
		cfg:       cfg,
		ec2Client: ec2.New(sess),
		ecsClient: ecs.New(sess),
	}, nil
}

func (ad *AWSDeployer) Deploy() error {
	logging.Info("Deploying to AWS ECS", zap.String("api", ad.cfg.GetClusterName()))

	taskDefinition, err := ad.registerTaskDefinition()
	if err != nil {
		return fmt.Errorf("failed to register task definition: %w", err)
	}

	err = ad.updateService(taskDefinition)
	if err != nil {
		return fmt.Errorf("failed to update service: %w", err)
	}

	logging.Info("Deployment to AWS ECS completed successfully")
	return nil
}

func (ad *AWSDeployer) registerTaskDefinition() (*string, error) {
	input := &ecs.RegisterTaskDefinitionInput{
		Family: aws.String(ad.cfg.GetClusterName()),
		ContainerDefinitions: []*ecs.ContainerDefinition{
			{
				Name:  aws.String(ad.cfg.GetClusterName()),
				Image: aws.String(fmt.Sprintf("%s:%s", ad.cfg.GetClusterName(), "latest")), // You may want to adjust this
				PortMappings: []*ecs.PortMapping{
					{
						ContainerPort: aws.Int64(int64(ad.cfg.GetServerPort())),
						HostPort:      aws.Int64(int64(ad.cfg.GetServerPort())),
					},
				},
			},
		},
	}

	result, err := ad.ecsClient.RegisterTaskDefinition(input)
	if err != nil {
		return nil, err
	}

	return result.TaskDefinition.TaskDefinitionArn, nil
}

func (ad *AWSDeployer) updateService(taskDefinitionArn *string) error {
	input := &ecs.UpdateServiceInput{
		Cluster:        aws.String(ad.cfg.GetClusterName()),
		Service:        aws.String(ad.cfg.GetClusterName()),
		TaskDefinition: taskDefinitionArn,
	}

	_, err := ad.ecsClient.UpdateService(input)
	return err
}
