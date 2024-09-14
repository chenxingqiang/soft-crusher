// File: internal/deployment/cloud_service_purchaser.go

package deployment

import (
	"fmt"
	"time"

	"github.com/chenxingqiang/soft-crusher/pkg/logging"
	"go.uber.org/zap"
)

// ... existing code ...

func (csp *CloudServicePurchaser) PurchaseCloudService() error {
	logging.Info("Starting cloud service purchase",
		zap.String("provider", csp.Config.CloudProvider),
		zap.String("clusterName", csp.Config.ClusterName))

	start := time.Now()
	var err error

	switch csp.Config.CloudProvider {
	case "aliyun":
		err = csp.purchaseAliyunACK()
	case "aws":
		err = csp.purchaseAWSECS()
	default:
		err = fmt.Errorf("unsupported cloud provider: %s", csp.Config.CloudProvider)
	}

	duration := time.Since(start)

	if err != nil {
		logging.Error("Failed to purchase cloud service",
			zap.Error(err),
			zap.Duration("duration", duration))
		return fmt.Errorf("failed to purchase cloud service: %w", err)
	}

	logging.Info("Successfully purchased cloud service",
		zap.Duration("duration", duration))
	return nil
}

func (csp *CloudServicePurchaser) purchaseAliyunACK() error {
	// ... existing Aliyun purchase code ...

	logging.Info("Aliyun ACK cluster purchase initiated",
		zap.String("clusterName", csp.Config.ClusterName),
		zap.Int("nodeCount", csp.Config.NodeCount))

	// Add error checking and more detailed logging
	if response.GetHttpStatus() != 200 {
		return fmt.Errorf("Aliyun API returned non-200 status: %d, content: %s",
			response.GetHttpStatus(), response.GetHttpContentString())
	}

	logging.Info("Aliyun ACK cluster purchase completed",
		zap.String("response", response.GetHttpContentString()))
	return nil
}

func (csp *CloudServicePurchaser) purchaseAWSECS() error {
	// ... existing AWS purchase code ...

	logging.Info("AWS ECS cluster purchase initiated",
		zap.String("clusterName", csp.Config.ClusterName))

	// Add error checking and more detailed logging
	if result.Cluster == nil || result.Cluster.ClusterArn == nil {
		return fmt.Errorf("AWS ECS cluster creation returned unexpected nil values")
	}

	logging.Info("AWS ECS cluster purchase completed",
		zap.String("clusterArn", *result.Cluster.ClusterArn))
	return nil
}
