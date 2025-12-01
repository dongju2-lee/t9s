#!/bin/bash

set -e

echo "🚀 Installing T9s - Terraform Infrastructure Manager"
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
echo "🔨 Building t9s..."
go build -o t9s ./cmd/t9s

# Check if /usr/local/bin is writable
if [ -w /usr/local/bin ]; then
    echo "📦 Installing t9s to /usr/local/bin..."
    mv t9s /usr/local/bin/
else
    echo "📦 Installing t9s to /usr/local/bin (requires sudo)..."
    sudo mv t9s /usr/local/bin/
fi

# Verify installation
if command -v t9s &> /dev/null; then
    echo ""
    echo "✅ T9s installed successfully!"
    echo ""
    echo "Run 't9s' to start the application"
    t9s --version
else
    echo "❌ Installation failed"
    exit 1
fi
