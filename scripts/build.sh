#!/bin/bash

set -e

echo "Building Soft-Crusher..."

# Build backend
go build -o soft-crusher cmd/soft-crusher/main.go

# Build frontend
cd frontend
npm install
npm run build
cd ..

echo "Build completed successfully."