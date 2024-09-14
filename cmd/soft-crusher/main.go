// File: cmd/soft-crusher/main.go

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chenxingqiang/soft-crusher/internal/auth"
	"github.com/chenxingqiang/soft-crusher/internal/config"
	"github.com/chenxingqiang/soft-crusher/internal/dashboard"
	"github.com/chenxingqiang/soft-crusher/internal/deployment"
	"github.com/chenxingqiang/soft-crusher/internal/plugin"
	"github.com/chenxingqiang/soft-crusher/pkg/logging"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	logging.SetLogger(logger)

	// Load configuration
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		logging.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Initialize components
	dashboard.Init()
	pluginManager := initializePluginManager(cfg)
	deployer := initializeDeployer(cfg)
	cloudServicePurchaser := deployment.NewCloudServicePurchaser(cfg)
	cloudDeploymentWorkflow := deployment.NewCloudDeploymentWorkflow(cfg)

	// Set up HTTP server
	mux := http.NewServeMux()
	setupRoutes(mux, cfg, deployer, cloudServicePurchaser, cloudDeploymentWorkflow, pluginManager)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Start server in a goroutine
	go func() {
		logging.Info("Starting server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logging.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logging.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logging.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logging.Info("Server exiting")
}

func initializePluginManager(cfg *config.Config) *plugin.PluginManager {
	pluginManager := plugin.NewPluginManager()
	if err := pluginManager.LoadPlugins(cfg.PluginDir); err != nil {
		logging.Error("Failed to load plugins", zap.Error(err))
	}
	return pluginManager
}

func initializeDeployer(cfg *config.Config) deployment.Deployer {
	var deployer deployment.Deployer
	var err error

	switch cfg.DeploymentTarget {
	case "kubernetes":
		deployer = deployment.NewKubernetesDeployer(cfg)
	case "aws":
		deployer, err = deployment.NewAWSDeployer(cfg)
	case "gcp":
		deployer, err = deployment.NewGCPDeployer(cfg)
	default:
		logging.Fatal("Unsupported deployment target", zap.String("target", cfg.DeploymentTarget))
	}

	if err != nil {
		logging.Fatal("Failed to create deployer", zap.Error(err))
	}

	return deployer
}

func setupRoutes(mux *http.ServeMux, cfg *config.Config, deployer deployment.Deployer,
	cloudServicePurchaser *deployment.CloudServicePurchaser,
	cloudDeploymentWorkflow *deployment.CloudDeploymentWorkflow,
	pluginManager *plugin.PluginManager) {

	mux.HandleFunc("/api/login", auth.Login)
	mux.HandleFunc("/api/dashboard", auth.JWTMiddleware(dashboard.EnhancedDashboardHandler))
	mux.HandleFunc("/api/deploy", auth.JWTMiddleware(handleDeploy(deployer)))
	mux.HandleFunc("/api/plugins", auth.JWTMiddleware(pluginManager.HandleListPlugins))
	mux.HandleFunc("/api/plugins/execute", auth.JWTMiddleware(pluginManager.HandleExecutePlugin))
	mux.HandleFunc("/api/purchase-cloud-service", auth.JWTMiddleware(handlePurchaseCloudService(cloudServicePurchaser)))
	mux.HandleFunc("/api/deploy-to-cloud", auth.JWTMiddleware(handleDeployToCloud(cfg, cloudDeploymentWorkflow)))

	// Serve the React frontend
	fs := http.FileServer(http.Dir("./frontend/build"))
	mux.Handle("/", fs)
}

func handleDeploy(deployer deployment.Deployer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := deployer.Deploy(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Deployment successful"))
	}
}

func handlePurchaseCloudService(purchaser *deployment.CloudServicePurchaser) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := purchaser.PurchaseCloudService(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Cloud service purchased successfully"))
	}
}

func handleDeployToCloud(cfg *config.Config, workflow *deployment.CloudDeploymentWorkflow) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var deployRequest struct {
			CloudProvider string `json:"cloudProvider"`
			ClusterName   string `json:"clusterName"`
			NodeCount     int    `json:"nodeCount"`
		}

		if err := json.NewDecoder(r.Body).Decode(&deployRequest); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Update config with request data
		cfg.CloudProvider = deployRequest.CloudProvider
		cfg.ClusterName = deployRequest.ClusterName
		cfg.NodeCount = deployRequest.NodeCount

		if err := workflow.Execute(); err != nil {
			logging.Error("Cloud deployment failed", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Cloud deployment completed successfully"))
	}
}
