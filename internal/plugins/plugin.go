// File: internal/plugin/plugin.go

package plugin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"plugin"
	"sync"

	"github.com/chenxingqiang/soft-crusher/pkg/logging"
	"go.uber.org/zap"
)

// Plugin represents a loaded plugin
type Plugin struct {
	Name    string
	Execute func() error
}

// PluginManager manages the loading and execution of plugins
type PluginManager struct {
	plugins map[string]Plugin
	mutex   sync.RWMutex
}

// NewPluginManager creates a new PluginManager
func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make(map[string]Plugin),
	}
}

// LoadPlugins loads all plugins from the specified directory
func (pm *PluginManager) LoadPlugins(pluginDir string) error {
	logging.Info("Loading plugins", zap.String("directory", pluginDir))

	files, err := ioutil.ReadDir(pluginDir)
	if err != nil {
		return fmt.Errorf("failed to read plugin directory: %w", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".so" {
			pluginPath := filepath.Join(pluginDir, file.Name())
			if err := pm.loadPlugin(pluginPath); err != nil {
				logging.Error("Failed to load plugin", zap.String("plugin", file.Name()), zap.Error(err))
			}
		}
	}

	logging.Info("Plugins loaded", zap.Int("count", len(pm.plugins)))
	return nil
}

func (pm *PluginManager) loadPlugin(path string) error {
	p, err := plugin.Open(path)
	if err != nil {
		return fmt.Errorf("could not open plugin: %w", err)
	}

	symPlugin, err := p.Lookup("Plugin")
	if err != nil {
		return fmt.Errorf("could not find Plugin symbol: %w", err)
	}

	plugin, ok := symPlugin.(Plugin)
	if !ok {
		return fmt.Errorf("unexpected type from symbol")
	}

	pm.mutex.Lock()
	pm.plugins[plugin.Name] = plugin
	pm.mutex.Unlock()

	logging.Info("Plugin loaded", zap.String("name", plugin.Name))
	return nil
}

// GetLoadedPlugins returns a list of all loaded plugins
func (pm *PluginManager) GetLoadedPlugins() []Plugin {
	pm.mutex.RLock()
	defer pm.mutex.RUnlock()

	plugins := make([]Plugin, 0, len(pm.plugins))
	for _, p := range pm.plugins {
		plugins = append(plugins, p)
	}
	return plugins
}

// ExecutePlugin executes the specified plugin
func (pm *PluginManager) ExecutePlugin(name string) error {
	pm.mutex.RLock()
	plugin, exists := pm.plugins[name]
	pm.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("plugin not found: %s", name)
	}

	return plugin.Execute()
}

// HandleListPlugins is an HTTP handler that lists all loaded plugins
func (pm *PluginManager) HandleListPlugins(w http.ResponseWriter, r *http.Request) {
	plugins := pm.GetLoadedPlugins()
	pluginNames := make([]string, len(plugins))
	for i, p := range plugins {
		pluginNames[i] = p.Name
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pluginNames)
}

// HandleExecutePlugin is an HTTP handler that executes a specified plugin
func (pm *PluginManager) HandleExecutePlugin(w http.ResponseWriter, r *http.Request) {
	pluginName := r.URL.Query().Get("name")
	if pluginName == "" {
		http.Error(w, "Plugin name is required", http.StatusBadRequest)
		return
	}

	err := pm.ExecutePlugin(pluginName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to execute plugin: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Plugin %s executed successfully", pluginName)))
}
