#!/bin/bash

# Build script for GoChatApp (production)
set -e

echo "Building GoChatApp for production..."

# Build React frontend
echo "Building React frontend..."
cd frontend
npm install
npm run build
cd ..

# Build Go backend
echo "Building Go backend..."
go build -o gochatapp .

echo ""
echo "✓ Build complete!"
echo ""
echo "Output:"
echo "  - Frontend: ./dist/"
echo "  - Backend:  ./gochatapp"
echo ""
echo "Run with: ./gochatapp"
