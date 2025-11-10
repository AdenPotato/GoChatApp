#!/bin/bash

# Start development servers for GoChatApp
echo "Starting GoChatApp development servers..."

# Function to cleanup background processes on exit
cleanup() {
    echo ""
    echo "Shutting down servers..."
    kill $BACKEND_PID $FRONTEND_PID 2>/dev/null
    exit
}

# Set trap to cleanup on script exit
trap cleanup SIGINT SIGTERM

# Start Go backend in background
echo "Starting Go backend server on :8080..."
go run main.go &
BACKEND_PID=$!

# Wait a bit for backend to start
sleep 2

# Start Vite frontend in background
echo "Starting Vite frontend server on :5173..."
cd frontend && npm run dev &
FRONTEND_PID=$!

echo ""
echo "✓ Backend running on http://localhost:8080"
echo "✓ Frontend running on http://localhost:5173"
echo ""
echo "Press Ctrl+C to stop both servers"

# Wait for both processes
wait
