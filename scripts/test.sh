#!/bin/bash

set -e

echo "Running tests for Soft-Crusher..."

# Run backend tests
go test ./...

# Run frontend tests
cd frontend
npm test
cd ..

echo "All tests completed successfully."