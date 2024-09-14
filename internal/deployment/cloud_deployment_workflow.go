// File: internal/deployment/cloud_deployment_workflow.go

package deployment

import (
	"fmt"
	"time"

	"github.com/chenxingqiang/soft-crusher/internal/config"
	"github.com/chenxingqiang/soft-crusher/pkg/logging"
	"go.uber.org/zap"
)

type CloudDeploymentWorkflow struct {
	Config           *config.Config
	Purchaser        *CloudServicePurchaser
	Deployer         Deployer
	ReadinessChecker *ClusterReadinessChecker
}

func NewCloudDeploymentWorkflow(cfg *config.Config) *CloudDeploymentWorkflow {
	return &CloudDeploymentWorkflow{
		Config:           cfg,
		Purchaser:        NewCloudServicePurchaser(cfg),
		Deployer:         NewDeployer(cfg),
		ReadinessChecker: NewClusterReadinessChecker(cfg),
	}
}

func (cdw *CloudDeploymentWorkflow) Execute() error {
	logging.Info("Starting cloud deployment workflow")

	// Step 1: Purchase cloud service
	err := cdw.Purchaser.PurchaseCloudService()
	if err != nil {
		return fmt.Errorf("failed to purchase cloud service: %w", err)
	}

	// Step 2: Wait for cluster to be ready
	err = cdw.waitForClusterReady()
	if err != nil {
		return fmt.Errorf("cluster did not become ready: %w", err)
	}

	// Step 3: Deploy application
	err = cdw.Deployer.Deploy()
	if err != nil {
		return fmt.Errorf("failed to deploy application: %w", err)
	}

	logging.Info("Cloud deployment workflow completed successfully")
	return nil
}

func (cdw *CloudDeploymentWorkflow) waitForClusterReady() error {
	logging.Info("Waiting for cluster to be ready")

	maxWaitTime := 15 * time.Minute
	checkInterval := 30 * time.Second
	startTime := time.Now()

	for time.Since(startTime) < maxWaitTime {
		ready, err := cdw.ReadinessChecker.CheckClusterReady()
		if err != nil {
			logging.Error("Error checking cluster readiness", zap.Error(err))
		} else if ready {
			logging.Info("Cluster is ready")
			return nil
		}

		logging.Info("Cluster not yet ready, waiting...",
			zap.Duration("elapsed", time.Since(startTime)))
		time.Sleep(checkInterval)
	}

	return fmt.Errorf("cluster did not become ready within %v", maxWaitTime)
}
