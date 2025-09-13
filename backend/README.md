# Spark Park Cricket Backend

[![Backend CI/CD](https://github.com/luffybhaagi/spark-park-cricket/actions/workflows/backend-ci.yml/badge.svg)](https://github.com/luffybhaagi/spark-park-cricket/actions/workflows/backend-ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.23-blue.svg)](https://golang.org/)
[![Go Report Card](https://goreportcard.com/badge/github.com/luffybhaagi/spark-park-cricket)](https://goreportcard.com/report/github.com/luffybhaagi/spark-park-cricket)
[![Coverage](https://codecov.io/gh/luffybhaagi/spark-park-cricket/branch/main/graph/badge.svg)](https://codecov.io/gh/luffybhaagi/spark-park-cricket)

A Go backend service with Supabase integration for the Spark Park Cricket application.

## Features

- HTTP server with health checks
- Supabase database integration
- Environment-based configuration
- Database health monitoring

## Setup

1. Copy the environment example file:
   ```bash
   cp env.example .env
   ```

2. Update `.env` with your Supabase credentials:
   ```
   SUPABASE_URL=your_supabase_project_url
   SUPABASE_API_KEY=your_supabase_anon_key
   PORT=8080
   ```

3. Install dependencies:
   ```bash
   go mod tidy
   ```

4. Run the server:
   ```bash
   go run main.go
   ```

## API Endpoints

- `GET /` - Welcome message
- `GET /health` - Basic health check
- `GET /db-health` - Database connection health check

## Project Structure

```
backend/
├── config/          # Configuration management
├── database/        # Database client and operations
├── main.go         # Application entry point
├── go.mod          # Go module file
├── env.example     # Environment variables template
└── README.md       # This file
```
