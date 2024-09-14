package config

import (
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	APIPort        int    `yaml:"api_port"`
	OutputDir      string `yaml:"output_dir"`
	EnableSwagger  bool   `yaml:"enable_swagger"`
	APIPrefix      string `yaml:"api_prefix"`
	ValidationType string `yaml:"validation_type"`

	CloudProvider string

	// Aliyun specific
	AliyunRegion          string
	AliyunAccessKeyID     string
	AliyunAccessKeySecret string

	// AWS specific
	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string

	// Common cluster configuration
	ClusterName        string
	KubernetesVersion  string
	WorkerInstanceType string
	NodeCount          int
}

func LoadConfig(filename string) (*Config, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	return &config, nil
}

func (c *Config) Validate() error {
	if c.APIPort <= 0 || c.APIPort > 65535 {
		return fmt.Errorf("invalid API port: %d", c.APIPort)
	}

	if c.OutputDir == "" {
		return fmt.Errorf("output directory is required")
	}

	if c.ValidationType != "basic" && c.ValidationType != "strict" {
		return fmt.Errorf("invalid validation type: %s", c.ValidationType)
	}

	return nil
}
