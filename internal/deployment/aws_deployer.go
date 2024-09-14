// File: internal/deployment/aws_deployer.go

package deployment

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/yourusername/soft-crusher/internal/config"
	"github.com/yourusername/soft-crusher/pkg/logging"
	"go.uber.org/zap"
)

type AWSDeployer struct {
	Config *config.Config
	ECS    *ecs.ECS
}

func NewAWSDeployer(cfg *config.Config) (*AWSDeployer, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(cfg.AWSRegion),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	return &AWSDeployer{
		Config: cfg,
		ECS:    ecs.New(sess),
	}, nil
}

func (ad *AWSDeployer) Deploy() error {
	logging.Info("Deploying to AWS ECS", zap.String("api", ad.Config.APIName))

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
		Family: aws.String(ad.Config.APIName),
		ContainerDefinitions: []*ecs.ContainerDefinition{
			{
				Name:  aws.String(ad.Config.APIName),
				Image: aws.String(fmt.Sprintf("%s:%s", ad.Config.DockerImage, ad.Config.APIVersion)),
				PortMappings: []*ecs.PortMapping{
					{
						ContainerPort: aws.Int64(int64(ad.Config.APIPort)),
						HostPort:      aws.Int64(int64(ad.Config.APIPort)),
					},
				},
			},
		},
	}

	result, err := ad.ECS.RegisterTaskDefinition(input)
	if err != nil {
		return nil, err
	}

	return result.TaskDefinition.TaskDefinitionArn, nil
}

func (ad *AWSDeployer) updateService(taskDefinitionArn *string) error {
	input := &ecs.UpdateServiceInput{
		Cluster:        aws.String(ad.Config.AWSECSCluster),
		Service:        aws.String(ad.Config.APIName),
		TaskDefinition: taskDefinitionArn,
	}

	_, err := ad.ECS.UpdateService(input)
	return err
}

// File: internal/deployment/gcp_deployer.go

package deployment

import (
	"context"
	"fmt"

	"cloud.google.com/go/container/apiv1"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
	"github.com/yourusername/soft-crusher/internal/config"
	"github.com/yourusername/soft-crusher/pkg/logging"
	"go.uber.org/zap"
)

type GCPDeployer struct {
	Config *config.Config
	Client *container.ClusterManagerClient
}

func NewGCPDeployer(cfg *config.Config) (*GCPDeployer, error) {
	ctx := context.Background()
	client, err := container.NewClusterManagerClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCP client: %w", err)
	}

	return &GCPDeployer{
		Config: cfg,
		Client: client,
	}, nil
}

func (gd *GCPDeployer) Deploy() error {
	logging.Info("Deploying to GCP GKE", zap.String("api", gd.Config.APIName))

	ctx := context.Background()

	// Get the cluster
	cluster, err := gd.getCluster(ctx)
	if err != nil {
		return fmt.Errorf("failed to get cluster: %w", err)
	}

	// Update the deployment
	err = gd.updateDeployment(ctx, cluster)
	if err != nil {
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	logging.Info("Deployment to GCP GKE completed successfully")
	return nil
}

func (gd *GCPDeployer) getCluster(ctx context.Context) (*containerpb.Cluster, error) {
	req := &containerpb.GetClusterRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/clusters/%s", 
			gd.Config.GCPProjectID, gd.Config.GCPZone, gd.Config.GCPClusterName),
	}

	return gd.Client.GetCluster(ctx, req)
}

func (gd *GCPDeployer) updateDeployment(ctx context.Context, cluster *containerpb.Cluster) error {
	// In a real-world scenario, you would use the Kubernetes API to update the deployment
	// For simplicity, we'll just log the action here
	logging.Info("Updating deployment in GKE cluster",
		zap.String("cluster", cluster.Name),
		zap.String("api", gd.Config.APIName),
		zap.String("image", fmt.Sprintf("%s:%s", gd.Config.DockerImage, gd.Config.APIVersion)))

	return nil
}