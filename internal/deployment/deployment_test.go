// File: internal/deployment/deployment_test.go

package deployment

import (
	"testing"

	"github.com/chenxingqiang/soft-crusher/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestNewDeployer(t *testing.T) {
	tests := []struct {
		name             string
		deploymentTarget string
		expectedType     interface{}
		expectError      bool
	}{
		{
			name:             "Kubernetes Deployer",
			deploymentTarget: "kubernetes",
			expectedType:     &KubernetesDeployer{},
			expectError:      false,
		},
		{
			name:             "AWS Deployer",
			deploymentTarget: "aws",
			expectedType:     &AWSDeployer{},
			expectError:      false,
		},
		{
			name:             "GCP Deployer",
			deploymentTarget: "gcp",
			expectedType:     &GCPDeployer{},
			expectError:      false,
		},
		{
			name:             "Unsupported Deployer",
			deploymentTarget: "unsupported",
			expectedType:     nil,
			expectError:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				DeploymentTarget: tt.deploymentTarget,
			}
			deployer, err := NewDeployer(cfg)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, deployer)
			} else {
				assert.NoError(t, err)
				assert.IsType(t, tt.expectedType, deployer)
			}
		})
	}
}
