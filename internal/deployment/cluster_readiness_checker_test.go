// File: internal/deployment/cluster_readiness_checker_test.go

package deployment

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yourusername/soft-crusher/internal/config"
)

func TestClusterReadinessChecker(t *testing.T) {
	tests := []struct {
		name          string
		cloudProvider string
		expectedError bool
	}{
		{
			name:          "Aliyun Provider",
			cloudProvider: "aliyun",
			expectedError: false,
		},
		{
			name:          "AWS Provider",
			cloudProvider: "aws",
			expectedError: false,
		},
		{
			name:          "Unsupported Provider",
			cloudProvider: "unsupported",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{
				CloudProvider: tt.cloudProvider,
				// Add other necessary config fields
			}

			checker := NewClusterReadinessChecker(cfg)
			_, err := checker.CheckClusterReady()

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
