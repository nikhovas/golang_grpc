FROM node:20-alpine

WORKDIR /app

# Copy package files
COPY package.json ./

# Install dependencies
RUN npm install

# Copy source code
COPY . .

# Expose port
EXPOSE 3000

# Start dev server (for development)
# For production, uncomment the build step and use "vite preview"
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]

