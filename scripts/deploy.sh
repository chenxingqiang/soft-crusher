#!/bin/bash

set -e

echo "Deploying Soft-Crusher..."

# Build Docker image
docker build -t soft-crusher:latest .

# Push to container registry (replace with your registry)
# docker push your-registry/soft-crusher:latest

# Deploy to Kubernetes (assumes kubectl is configured)
kubectl apply -f k8s/

echo "Deployment completed successfully."