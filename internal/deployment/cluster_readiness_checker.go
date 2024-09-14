// File: internal/deployment/cluster_readiness_checker.go

package deployment

import (
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/chenxingqiang/soft-crusher/internal/config"
	"github.com/chenxingqiang/soft-crusher/pkg/logging"
	"go.uber.org/zap"
)

type ClusterReadinessChecker struct {
	Config *config.Config
}

func NewClusterReadinessChecker(cfg *config.Config) *ClusterReadinessChecker {
	return &ClusterReadinessChecker{Config: cfg}
}

func (crc *ClusterReadinessChecker) CheckClusterReady() (bool, error) {
	switch crc.Config.CloudProvider {
	case "aliyun":
		return crc.checkAliyunClusterReady()
	case "aws":
		return crc.checkAWSClusterReady()
	default:
		return false, fmt.Errorf("unsupported cloud provider: %s", crc.Config.CloudProvider)
	}
}

func (crc *ClusterReadinessChecker) checkAliyunClusterReady() (bool, error) {
	client, err := sdk.NewClientWithAccessKey(
		crc.Config.AliyunRegion,
		crc.Config.AliyunAccessKeyID,
		crc.Config.AliyunAccessKeySecret,
	)
	if err != nil {
		return false, fmt.Errorf("failed to create Aliyun client: %w", err)
	}

	request := requests.NewCommonRequest()
	request.Method = "GET"
	request.Scheme = "https"
	request.Domain = "cs.aliyuncs.com"
	request.Version = "2015-12-15"
	request.ApiName = "DescribeClusterDetail"
	request.QueryParams["RegionId"] = crc.Config.AliyunRegion
	request.QueryParams["ClusterName"] = crc.Config.ClusterName

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		return false, fmt.Errorf("failed to check Aliyun cluster status: %w", err)
	}

	if response.GetHttpStatus() != 200 {
		return false, fmt.Errorf("Aliyun API returned non-200 status: %d", response.GetHttpStatus())
	}

	// Parse the response and check if the cluster state is "running"
	// This is a simplified example. You should parse the actual JSON response.
	if response.GetHttpContentString() != "" {
		logging.Info("Aliyun cluster is ready")
		return true, nil
	}

	logging.Info("Aliyun cluster is not yet ready")
	return false, nil
}

func (crc *ClusterReadinessChecker) checkAWSClusterReady() (bool, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(crc.Config.AWSRegion),
	})
	if err != nil {
		return false, fmt.Errorf("failed to create AWS session: %w", err)
	}

	svc := ecs.New(sess)

	input := &ecs.DescribeClustersInput{
		Clusters: []*string{
			aws.String(crc.Config.ClusterName),
		},
	}

	result, err := svc.DescribeClusters(input)
	if err != nil {
		return false, fmt.Errorf("failed to describe AWS ECS cluster: %w", err)
	}

	if len(result.Clusters) == 0 {
		return false, fmt.Errorf("AWS ECS cluster not found")
	}

	cluster := result.Clusters[0]
	if *cluster.Status == "ACTIVE" {
		logging.Info("AWS ECS cluster is ready")
		return true, nil
	}

	logging.Info("AWS ECS cluster is not yet ready", zap.String("status", *cluster.Status))
	return false, nil
}
