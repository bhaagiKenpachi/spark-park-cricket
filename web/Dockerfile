# Build stage
FROM node:20-alpine AS builder

# Set working directory
WORKDIR /app

# Copy package files
COPY package*.json ./

# Install all dependencies
RUN npm ci

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Production stage
FROM node:20-alpine

# Set working directory
WORKDIR /app

# Install curl for health checks
RUN apk add --no-cache curl

# Copy package files
COPY package*.json ./

# Install production dependencies + TypeScript (needed for next.config.ts at runtime)
RUN npm ci --omit=dev --ignore-scripts || npm ci --production --ignore-scripts

# Install TypeScript separately (needed for next.config.ts)
RUN npm install --save-prod typescript @types/node

# Copy built application from builder
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/next.config.* ./

# Expose port
EXPOSE 3000

# Health check with longer start period
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=5 \
    CMD curl -f http://localhost:3000 || exit 1

# Set environment variables
ENV NODE_ENV=production
ENV PORT=3000

# Run the application in production mode
CMD ["npm", "start"]