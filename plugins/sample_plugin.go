// File: plugins/sample_plugin.go

package main

import (
	"fmt"

	"github.com/yourusername/soft-crusher/pkg/logging"
	"go.uber.org/zap"
)

type SamplePlugin struct{}

func (p *SamplePlugin) Info() PluginInfo {
	return PluginInfo{
		Name:        "SamplePlugin",
		Description: "A sample plugin for demonstration purposes",
		Version:     "1.0.0",
	}
}

func (p *SamplePlugin) Execute() error {
	logging.Info("Executing sample plugin")
	fmt.Println("Hello from the sample plugin!")
	return nil
}

// This is the symbol that the plugin manager will look for
var Plugin SamplePlugin
