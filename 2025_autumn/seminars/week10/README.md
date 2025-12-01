# Week 10: gRPC-Web with Envoy Proxy

This project demonstrates a gRPC server with Envoy proxy for gRPC-Web support and a React TypeScript web client.

## Architecture

```
React Web Client (Port 3000)
    ↓ HTTP/gRPC-Web
Envoy Proxy (Port 8080)
    ↓ gRPC
gRPC Server (Port 50051)
    ↓
Redis (Port 6379)
```

## Prerequisites

- Docker and Docker Compose
- Go 1.22+ (for building the server)
- Node.js 18+ (for the web client)
- protoc (Protocol Buffers compiler)
- protoc-gen-grpc-web plugin

## Quick Start

### 1. Start the Backend Services

```bash
cd server
docker-compose up -d
```

This will start:
- Redis (data storage)
- gRPC Server (port 50051)
- Envoy Proxy (port 8080)

### 2. Set Up the Web Client

```bash
cd web

# Install dependencies
npm install

# Install protoc-gen-grpc-web globally (if not already installed)
npm install -g protoc-gen-grpc-web

# Generate gRPC-Web client code from proto files
npm run generate

# Start the development server
npm run dev
```

The web client will be available at http://localhost:3000

## Project Structure

```
week10/
├── api/                    # Protocol Buffer definitions
│   └── server.proto
├── server/                 # gRPC Server (Go)
│   ├── main.go
│   ├── Dockerfile
│   └── docker-compose.yml
├── envoy/                  # Envoy Proxy configuration
│   └── envoy.yaml
└── web/                    # React TypeScript Web Client
    ├── src/
    │   ├── App.tsx
    │   └── ...
    └── package.json
```

## Services

### gRPC Server
- **Port**: 50051
- **Language**: Go
- **Service**: KeyValueService
  - `SetValue(key, value)` - Store a key-value pair
  - `GetValue(key)` - Retrieve a value by key

### Envoy Proxy
- **Port**: 8080 (HTTP/gRPC-Web)
- **Admin Port**: 9901
- **Function**: Translates gRPC-Web requests to gRPC

### Web Client
- **Port**: 3000
- **Framework**: React + TypeScript + Vite
- **Protocol**: gRPC-Web

## Testing

1. Open http://localhost:3000 in your browser
2. Enter a key and value, then click "Set Value"
3. Enter a key (leave value empty), then click "Get Value"

## Troubleshooting

### Envoy not connecting to gRPC server
- Ensure both services are on the same Docker network
- Check that the gRPC server is running: `docker-compose ps`

### Web client can't connect to Envoy
- Verify Envoy is running on port 8080
- Check browser console for CORS errors
- Ensure `VITE_ENVOY_URL` in `.env` matches your Envoy URL

### Proto generation fails
- Install protoc: https://grpc.io/docs/protoc-installation/
- Install protoc-gen-grpc-web: `npm install -g protoc-gen-grpc-web`
- Verify protoc-gen-grpc-web is in your PATH

## Development

### Regenerating Proto Files

If you modify `api/server.proto`:

**Server (Go):**
```bash
cd tools
./proto-gen.sh
```

**Web Client (TypeScript):**
```bash
cd web
npm run generate
```

### Viewing Envoy Admin Interface

Visit http://localhost:9901 for Envoy's admin interface.

