# scripts/setup.sh - Initial setup script
#!/bin/bash

set -e

echo "Setting up TCR Game project..."

# Create data directory structure
mkdir -p data/players
touch data/players/.gitkeep

# Create log directory
mkdir -p logs

# Create necessary directories
mkdir -p web/static/css
mkdir -p web/static/js
mkdir -p web/static/images
mkdir -p web/templates

# Initialize Go module if not exists
if [ ! -f "go.mod" ]; then
    echo "Initializing Go module..."
    go mod init tcr-game
fi

# Download dependencies
echo "Downloading dependencies..."
go mod tidy

# Create .env file if not exists
if [ ! -f ".env" ]; then
    echo "Creating .env file..."
    cat > .env << EOF
# TCR Game Environment Variables
SERVER_PORT=8080
LOG_LEVEL=info
EOF
fi

echo "Setup complete! Run 'go run main.go' to start the server"