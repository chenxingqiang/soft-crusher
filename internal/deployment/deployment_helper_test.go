package deployment

import (
	"os"
	"testing"
"github.com/stretchr/testify/assert"
)

func TestDeploymentHelper(t *testing.T) {
	helper := NewDeploymentHelper("test-api", "1.0.0", 8080)

	t.Run("GenerateDockerfile", func(t *testing.T) {
		err := helper.GenerateDockerfile()
		assert.NoError(t, err)
		assert.FileExists(t, "Dockerfile")
		
		content, err := os.ReadFile("Dockerfile")
		assert.NoError(t, err)
		assert.Contains(t, string(content), "EXPOSE 8080")
		assert.Contains(t, string(content), "CMD [\"./soft-crusher-api\"]")

		os.Remove("Dockerfile")
	})

	t.Run("GenerateKubernetesManifests", func(t *testing.T) {
		err := helper.GenerateKubernetesManifests()
		assert.NoError(t, err)
		assert.FileExists(t, "kubernetes-manifests.yaml")
		
		content, err := os.ReadFile("kubernetes-manifests.yaml")
		assert.NoError(t, err)
		assert.Contains(t, string(content), "name: test-api")
		assert.Contains(t, string(content), "image: test-api:1.0.0")
		assert.Contains(t, string(content), "containerPort: 8080")

		os.Remove("kubernetes-manifests.yaml")
	})

	t.Run("GenerateDockerComposeFile", func(t *testing.T) {
		err := helper.GenerateDockerComposeFile()
		assert.NoError(t, err)
		assert.FileExists(t, "docker-compose.yaml")
		
		content, err := os.ReadFile("docker-compose.yaml")
		assert.NoError(t, err)
		assert.Contains(t, string(content), "test-api:")
		assert.Contains(t, string(content), "- \"8080:8080\"")

		os.Remove("docker-compose.yaml")
	})
}