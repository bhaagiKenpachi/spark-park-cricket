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

# Set build-time environment variables for Next.js
ARG NEXT_PUBLIC_API_URL=https://spark-park-dev.dojima.foundation/api/v1
ARG NEXT_PUBLIC_GRAPHQL_URL=https://spark-park-dev.dojima.foundation/api/v1/graphql
ARG NEXT_PUBLIC_ADSENSE_CLIENT_ID=ca-pub-5474524579770573
ARG NEXT_PUBLIC_ENABLE_ADS
ARG NEXT_PUBLIC_SSE_BALL_EVENT_TIMEOUT
ARG NEXT_PUBLIC_SSE_INACTIVITY_WARNING_THRESHOLD
ARG NEXT_PUBLIC_SSE_INACTIVITY_DISCONNECT_THRESHOLD
ARG NEXT_PUBLIC_POSTHOG_KEY
ARG NEXT_PUBLIC_POSTHOG_HOST
ARG NEXT_PUBLIC_POSTHOG_DEBUG
ARG NODE_ENV

ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_GRAPHQL_URL=$NEXT_PUBLIC_GRAPHQL_URL
ENV NEXT_PUBLIC_ADSENSE_CLIENT_ID=$NEXT_PUBLIC_ADSENSE_CLIENT_ID
ENV NEXT_PUBLIC_ENABLE_ADS=$NEXT_PUBLIC_ENABLE_ADS
ENV NEXT_PUBLIC_SSE_BALL_EVENT_TIMEOUT=$NEXT_PUBLIC_SSE_BALL_EVENT_TIMEOUT
ENV NEXT_PUBLIC_SSE_INACTIVITY_WARNING_THRESHOLD=$NEXT_PUBLIC_SSE_INACTIVITY_WARNING_THRESHOLD
ENV NEXT_PUBLIC_SSE_INACTIVITY_DISCONNECT_THRESHOLD=$NEXT_PUBLIC_SSE_INACTIVITY_DISCONNECT_THRESHOLD
ENV NEXT_PUBLIC_POSTHOG_KEY=$NEXT_PUBLIC_POSTHOG_KEY
ENV NEXT_PUBLIC_POSTHOG_HOST=$NEXT_PUBLIC_POSTHOG_HOST
ENV NEXT_PUBLIC_POSTHOG_DEBUG=$NEXT_PUBLIC_POSTHOG_DEBUG
ENV NODE_ENV=$NODE_ENV

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