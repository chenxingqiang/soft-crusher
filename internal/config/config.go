package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Config represents the application configuration
type Config struct {
	Server struct {
		Port  int  `yaml:"port"`
		Debug bool `yaml:"debug"`
	} `yaml:"server"`

	Logging struct {
		Level       string   `yaml:"level"`
		File        string   `yaml:"file"`
		OutputPaths []string `yaml:"output_paths"`
	} `yaml:"logging"`

	Database struct {
		Type string `yaml:"type"`
		URI  string `yaml:"uri"`
		Name string `yaml:"name"`
	} `yaml:"database"`

	Auth struct {
		JWTSecret       string `yaml:"jwt_secret"`
		TokenExpiration string `yaml:"token_expiration"`
	} `yaml:"auth"`

	Deployment struct {
		Target    string `yaml:"target"`
		PluginDir string `yaml:"plugin_dir"`
	} `yaml:"deployment"`

	Cloud struct {
		Provider    string `yaml:"provider"`
		Region      string `yaml:"region"`
		ClusterName string `yaml:"cluster_name"`
		NodeCount   int    `yaml:"node_count"`
	} `yaml:"cloud"`

	AWS struct {
		AccessKeyID     string `yaml:"access_key_id"`
		SecretAccessKey string `yaml:"secret_access_key"`
	} `yaml:"aws"`

	GCP struct {
		ProjectID       string `yaml:"project_id"`
		CredentialsFile string `yaml:"credentials_file"`
	} `yaml:"gcp"`

	Kubernetes struct {
		Namespace      string `yaml:"namespace"`
		ServiceAccount string `yaml:"service_account"`
	} `yaml:"kubernetes"`

	Plugins struct {
		Enabled   bool   `yaml:"enabled"`
		Directory string `yaml:"directory"`
	} `yaml:"plugins"`

	Frontend struct {
		BuildDir string `yaml:"build_dir"`
	} `yaml:"frontend"`
}

// LoadConfig reads the config file and returns a Config struct
func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}

	// Read the config file
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// Unmarshal the YAML into the Config struct
	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return config, nil
}

// IsDebugMode returns true if debug mode is enabled
func (c *Config) IsDebugMode() bool {
	return c.Server.Debug
}

// GetServerPort returns the server port
func (c *Config) GetServerPort() int {
	return c.Server.Port
}

// GetLogLevel returns the log level
func (c *Config) GetLogLevel() string {
	return c.Logging.Level
}

// GetLogFile returns the log file path
func (c *Config) GetLogFile() string {
	return c.Logging.File
}

// GetOutputPaths returns the logging output paths
func (c *Config) GetOutputPaths() []string {
	return c.Logging.OutputPaths
}

// GetDatabaseURI returns the database URI
func (c *Config) GetDatabaseURI() string {
	return c.Database.URI
}

// GetDatabaseName returns the database name
func (c *Config) GetDatabaseName() string {
	return c.Database.Name
}

// GetJWTSecret returns the JWT secret
func (c *Config) GetJWTSecret() string {
	return c.Auth.JWTSecret
}

// GetTokenExpiration returns the token expiration time
func (c *Config) GetTokenExpiration() string {
	return c.Auth.TokenExpiration
}

// GetDeploymentTarget returns the deployment target
func (c *Config) GetDeploymentTarget() string {
	return c.Deployment.Target
}

// GetPluginDir returns the plugin directory
func (c *Config) GetPluginDir() string {
	return c.Deployment.PluginDir
}

// GetCloudProvider returns the cloud provider
func (c *Config) GetCloudProvider() string {
	return c.Cloud.Provider
}

// GetCloudRegion returns the cloud region
func (c *Config) GetCloudRegion() string {
	return c.Cloud.Region
}

// GetClusterName returns the cluster name
func (c *Config) GetClusterName() string {
	return c.Cloud.ClusterName
}

// GetNodeCount returns the node count
func (c *Config) GetNodeCount() int {
	return c.Cloud.NodeCount
}

// GetAWSAccessKeyID returns the AWS access key ID
func (c *Config) GetAWSAccessKeyID() string {
	return c.AWS.AccessKeyID
}

// GetAWSSecretAccessKey returns the AWS secret access key
func (c *Config) GetAWSSecretAccessKey() string {
	return c.AWS.SecretAccessKey
}

// GetGCPProjectID returns the GCP project ID
func (c *Config) GetGCPProjectID() string {
	return c.GCP.ProjectID
}

// GetGCPCredentialsFile returns the path to the GCP credentials file
func (c *Config) GetGCPCredentialsFile() string {
	return c.GCP.CredentialsFile
}

// GetKubernetesNamespace returns the Kubernetes namespace
func (c *Config) GetKubernetesNamespace() string {
	return c.Kubernetes.Namespace
}

// GetKubernetesServiceAccount returns the Kubernetes service account
func (c *Config) GetKubernetesServiceAccount() string {
	return c.Kubernetes.ServiceAccount
}

// IsPluginsEnabled returns true if plugins are enabled
func (c *Config) IsPluginsEnabled() bool {
	return c.Plugins.Enabled
}

// GetPluginsDirectory returns the plugins directory
func (c *Config) GetPluginsDirectory() string {
	return c.Plugins.Directory
}

// GetFrontendBuildDir returns the frontend build directory
func (c *Config) GetFrontendBuildDir() string {
	return c.Frontend.BuildDir
}
