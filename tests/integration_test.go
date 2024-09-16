// File: tests/integration_test.go
package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	baseURL        = "http://localhost:8080"
	defaultTimeout = 10 * time.Minute
)

var (
	apiKey  = os.Getenv("SOFT_CRUSHER_API_KEY")
	timeout = getEnvDuration("TEST_TIMEOUT", defaultTimeout)
)

func TestSoftCrusherWorkflow(t *testing.T) {
	require.NotEmpty(t, apiKey, "SOFT_CRUSHER_API_KEY environment variable is not set")

	client := &http.Client{Timeout: timeout}

	t.Run("AnalyzeCode", func(t *testing.T) {
		payload := map[string]string{"repoUrl": "https://github.com/maxence-charriere/go-app.git"}
		resp, err := sendRequest(client, "POST", "/api/analyze", payload)
		require.NoError(t, err, "Failed to send analyze request")
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Unexpected status code from analyze endpoint")

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err, "Failed to decode analyze response")
		assert.Contains(t, result, "functions", "Response does not contain 'functions' key")
		functions, ok := result["functions"].([]interface{})
		require.True(t, ok, "Functions is not a slice")
		assert.NotEmpty(t, functions, "No functions were analyzed")
	})

	t.Run("GenerateAPI", func(t *testing.T) {
		resp, err := sendRequest(client, "POST", "/api/generate", nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)
		assert.Contains(t, result, "apiSpec")
		assert.Contains(t, result["apiSpec"].(map[string]interface{}), "endpoints")
	})

	t.Run("EstimateCloudCost", func(t *testing.T) {
		payload := map[string]interface{}{
			"cloudProvider": "aws",
			"clusterName":   "test-cluster",
			"nodeCount":     3,
		}
		resp, err := sendRequest(client, "POST", "/api/estimate-cost", payload)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)
		assert.Contains(t, result, "estimatedCost")
		assert.Greater(t, result["estimatedCost"].(float64), 0.0)
	})

	t.Run("DeployToCloud", func(t *testing.T) {
		payload := map[string]interface{}{
			"cloudProvider": "aws",
			"clusterName":   "test-cluster",
			"nodeCount":     3,
		}
		resp, err := sendRequest(client, "POST", "/api/deploy-to-cloud", payload)
		require.NoError(t, err, "Failed to send deploy request")
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Unexpected status code from deploy endpoint")

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		require.NoError(t, err, "Failed to decode deploy response")
		assert.Contains(t, result, "status", "Response does not contain 'status' key")
		assert.Equal(t, "started", result["status"], "Deployment status is not 'started'")

		err = pollDeploymentStatus(client)
		require.NoError(t, err, "Deployment failed or timed out")
	})

	t.Run("GetDeployedAPIInfo", func(t *testing.T) {
		resp, err := sendRequest(client, "GET", "/api/deployed-api-info", nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]string
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)
		assert.Contains(t, result, "apiUrl")
		assert.Contains(t, result, "swaggerUrl")
	})

	t.Cleanup(func() {
		// Add cleanup logic here, e.g., deleting deployed resources
		_, err := sendRequest(client, "POST", "/api/cleanup", nil)
		if err != nil {
			t.Logf("Failed to clean up resources: %v", err)
		}
	})
}

func pollDeploymentStatus(client *http.Client) error {
	startTime := time.Now()
	for time.Since(startTime) < timeout {
		statusResp, err := sendRequest(client, "GET", "/api/deployment-status", nil)
		if err != nil {
			return fmt.Errorf("failed to get deployment status: %w", err)
		}

		var statusResult map[string]string
		err = json.NewDecoder(statusResp.Body).Decode(&statusResult)
		if err != nil {
			return fmt.Errorf("failed to decode deployment status: %w", err)
		}

		switch statusResult["status"] {
		case "completed":
			return nil
		case "failed":
			return fmt.Errorf("deployment failed")
		}

		time.Sleep(30 * time.Second)
	}
	return fmt.Errorf("deployment timed out")
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
func sendRequest(client *http.Client, method, path string, payload interface{}) (*http.Response, error) {
	var body bytes.Buffer
	if payload != nil {
		if err := json.NewEncoder(&body).Encode(payload); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, baseURL+path, &body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	return client.Do(req)
}
