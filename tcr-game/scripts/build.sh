# scripts/build.sh - Build script
#!/bin/bash

set -e

echo "Building TCR Game..."

# Create build directory
mkdir -p dist

# Build server
echo "Building server..."
go build -o dist/tcr-server ./main.go

# Build client (if exists)
if [ -f "cmd/client/main.go" ]; then
    echo "Building client..."
    go build -o dist/tcr-client ./cmd/client/main.go
fi

# Copy configuration files
echo "Copying config files..."
cp -r config dist/
cp -r data dist/
cp -r web dist/

echo "Build complete! Binary is in dist/tcr-server"



