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
)

const (
	baseURL = "http://localhost:8080"
	timeout = 10 * time.Minute
)

var apiKey = os.Getenv("SOFT_CRUSHER_API_KEY")

func TestSoftCrusherWorkflow(t *testing.T) {
	client := &http.Client{Timeout: timeout}

	t.Run("AnalyzeCode", func(t *testing.T) {
		payload := map[string]string{"repoUrl": "https://github.com/example/repo.git"}
		resp, err := sendRequest(client, "POST", "/api/analyze", payload)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)
		assert.Contains(t, result, "functions")
		assert.Greater(t, len(result["functions"].([]interface{})), 0)
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
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)
		assert.Contains(t, result, "status")
		assert.Equal(t, "started", result["status"])

		// Poll for deployment status
		startTime := time.Now()
		for time.Since(startTime) < timeout {
			statusResp, err := sendRequest(client, "GET", "/api/deployment-status", nil)
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, statusResp.StatusCode)

			var statusResult map[string]string
			err = json.NewDecoder(statusResp.Body).Decode(&statusResult)
			assert.NoError(t, err)

			if statusResult["status"] == "completed" {
				break
			} else if statusResult["status"] == "failed" {
				t.Fatal("Deployment failed")
			}

			time.Sleep(30 * time.Second)
		}

		if time.Since(startTime) >= timeout {
			t.Fatal("Deployment timed out")
		}
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
