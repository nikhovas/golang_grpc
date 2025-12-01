# gRPC-Web Client

React TypeScript web client for the Key-Value gRPC service.

## Prerequisites

- Node.js (v18 or higher)
- npm or yarn
- protoc (Protocol Buffers compiler)
- protoc-gen-grpc-web plugin

## Installation

1. Install protoc-gen-grpc-web:
   ```bash
   npm install -g protoc-gen-grpc-web
   ```
   
   Or download from: https://github.com/grpc/grpc-web/releases

2. Install dependencies:
   ```bash
   npm install
   ```

3. Generate gRPC-Web client code:
   ```bash
   npm run generate
   ```

4. Start the development server:
   ```bash
   npm run dev
   ```

The app will be available at http://localhost:3000

## Environment Variables

Create a `.env` file (or copy from `.env.example`):
```
VITE_ENVOY_URL=http://localhost:8080
```

## Building for Production

```bash
npm run build
```

The built files will be in the `dist` directory.

