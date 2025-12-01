#!/bin/bash

set -e

echo "🚀 Installing P9s - Terraform Infrastructure Manager"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.20 or higher."
    echo "   Visit: https://golang.org/dl/"
    exit 1
fi

# Get Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "✓ Found Go version: $GO_VERSION"

# Build the application
echo "🔨 Building p9s..."
go build -o p9s ./cmd/p9s

# Check if /usr/local/bin is writable
if [ -w /usr/local/bin ]; then
    echo "📦 Installing p9s to /usr/local/bin..."
    mv p9s /usr/local/bin/
else
    echo "📦 Installing p9s to /usr/local/bin (requires sudo)..."
    sudo mv p9s /usr/local/bin/
fi

# Verify installation
if command -v p9s &> /dev/null; then
    echo ""
    echo "✅ P9s installed successfully!"
    echo ""
    echo "Run 'p9s' to start the application"
    p9s --version
else
    echo "❌ Installation failed"
    exit 1
fi
