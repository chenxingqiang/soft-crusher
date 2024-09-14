// File: cmd/soft-crusher/main.go

package main

import (
	"github.com/yourusername/soft-crusher/internal/auth"
	"github.com/yourusername/soft-crusher/internal/config"
	"github.com/yourusername/soft-crusher/internal/dashboard"
	"github.com/yourusername/soft-crusher/internal/deployment"
	"github.com/yourusername/soft-crusher/internal/plugin"
	"github.com/yourusername/soft-crusher/pkg/logging"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		logging.Fatal("Failed to load configuration", zap.Error(err))
	}

	pluginManager := plugin.NewPluginManager()
	err = pluginManager.LoadPlugins(cfg.PluginDir)
	if err != nil {
		logging.Error("Failed to load plugins", zap.Error(err))
	}

	var deployer deployment.Deployer
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

	cloudServicePurchaser := deployment.NewCloudServicePurchaser(cfg)

	cloudDeploymentWorkflow := deployment.NewCloudDeploymentWorkflow(cfg)

	http.HandleFunc("/api/purchase-cloud-service", auth.JWTMiddleware(func(w http.ResponseWriter, r *http.Request) {
		err := cloudServicePurchaser.PurchaseCloudService()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Cloud service purchased successfully"))
	}))

	http.HandleFunc("/api/deploy-to-cloud", auth.JWTMiddleware(func(w http.ResponseWriter, r *http.Request) {
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

		err := cloudDeploymentWorkflow.Execute()
		if err != nil {
			logging.Error("Cloud deployment failed", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Cloud deployment completed successfully"))
	}))

	http.HandleFunc("/api/login", auth.Login)
	http.HandleFunc("/api/dashboard", auth.JWTMiddleware(dashboard.EnhancedDashboardHandler))
	http.HandleFunc("/api/deploy", auth.JWTMiddleware(func(w http.ResponseWriter, r *http.Request) {
		err := deployer.Deploy()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Deployment successful"))
	}))
	http.HandleFunc("/api/plugins", auth.JWTMiddleware(pluginManager.HandleListPlugins))
	http.HandleFunc("/api/plugins/execute", auth.JWTMiddleware(pluginManager.HandleExecutePlugin))

	// Serve the React frontend
	fs := http.FileServer(http.Dir("./frontend/build"))
	http.Handle("/", fs)

	logging.Info("Starting server on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		logging.Fatal("Failed to start server", zap.Error(err))
	}
}
