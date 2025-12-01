#!/bin/bash

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# Get the week10 directory (parent of tools/)
WEEK10_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

export PATH="${PATH}:${HOME}/go/bin"

# Check for protoc-gen-grpc-web (required for web client generation)
PROTOC_GEN_GRPC_WEB_PATH=$(which protoc-gen-grpc-web)

if [ -z "$PROTOC_GEN_GRPC_WEB_PATH" ]; then
    echo "Warning: protoc-gen-grpc-web not found. Web client files will not be generated."
    echo "Install it via: npm install -g protoc-gen-grpc-web"
    echo "Or download from: https://github.com/grpc/grpc-web/releases"
    GENERATE_WEB=false
else
    GENERATE_WEB=true
fi

# Change to week10 directory to ensure consistent paths
cd "$WEEK10_DIR"

# Generate Go files for server
echo "Generating Go files for server..."
mkdir -p server/api
protoc -I api --go_out server/api --go_opt paths=source_relative --go-grpc_out server/api --go-grpc_opt paths=source_relative \
api/*.proto

if [ $? -eq 0 ]; then
    echo "✓ Go files generated successfully in server/api/"
else
    echo "✗ Failed to generate Go files"
    exit 1
fi

# Generate JavaScript/TypeScript files for web client
if [ "$GENERATE_WEB" = true ]; then
    echo "Generating JavaScript/TypeScript files for web client..."
    mkdir -p web/src/generated
    
    # Check if ts-protoc-gen is available
    TS_PROTOC_GEN_PATH="$WEEK10_DIR/web/node_modules/.bin/protoc-gen-ts"
    
    if [ -f "$TS_PROTOC_GEN_PATH" ]; then
        echo "Using ts-protoc-gen for TypeScript message generation..."
        # Generate TypeScript messages
        protoc \
          -I=api \
          --plugin=protoc-gen-ts="$TS_PROTOC_GEN_PATH" \
          --ts_out=web/src/generated \
          api/*.proto
    else
        echo "ts-protoc-gen not found, using standard JS generation..."
        # Fallback: Generate JS with Google Closure compatibility
        protoc \
          -I=api \
          --js_out=import_style=commonjs,binary:web/src/generated \
          api/*.proto
    fi
    
    # Generate gRPC-Web client in TypeScript
    protoc \
      -I=api \
      --grpc-web_out=import_style=typescript,mode=grpcwebtext:web/src/generated \
      api/*.proto
    
    if [ $? -eq 0 ]; then
        echo "✓ Web client files generated successfully in web/src/generated/"
    else
        echo "✗ Failed to generate web client files"
        exit 1
    fi
fi

echo "Proto generation complete!"