import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
// @ts-ignore
import commonjs from 'vite-plugin-commonjs'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    react(),
    commonjs({
      filter: (id) => {
        // Transform generated files that use require()
        if (id.includes('generated') || id.includes('server_grpc_web_pb') || id.includes('server_pb')) {
          return true
        }
        return false
      }
    })
  ],
  optimizeDeps: {
    include: ['google-protobuf', 'grpc-web'],
    esbuildOptions: {
      mainFields: ['module', 'main']
    },
    // Force pre-bundling of generated files
    exclude: []
  },
  build: {
    commonjsOptions: {
      include: [/generated/, /node_modules/],
      transformMixedEsModules: true,
      // Transform require() calls
      requireReturnsDefault: 'auto'
    }
  },
  resolve: {
    dedupe: ['google-protobuf', 'grpc-web']
  },
  // Ensure CommonJS transformation for generated files
  ssr: {
    noExternal: ['google-protobuf', 'grpc-web']
  },
  server: {
    port: 3000,
    host: '0.0.0.0',
    proxy: {
      '/api.KeyValueService': {
        target: process.env.VITE_ENVOY_URL || 'http://localhost:8080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/api/, '')
      }
    }
  }
})

