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

## Performance Benchmarks

### Database Operations Performance (January 2025)

| Operation Type | Response Time | Throughput | Data Transfer | Status |
|----------------|---------------|------------|---------------|---------|
| **Series Operations** |
| Get Series | 55.01ms | 22 ops/sec | 242 bytes | ✅ Excellent |
| Create Series | 58.18ms | 18 ops/sec | 240 bytes | ✅ Excellent |
| Series with Matches | 65.73ms | 19.8 ops/sec | 242 bytes | ✅ Excellent |
| **Match Operations** |
| Create Match | 170.5ms | 5.6 ops/sec | 381 bytes | ⚠️ Good |
| Update Match | 115.6ms | 9.8 ops/sec | 381 bytes | ⚠️ Good |
| **Scorecard Operations** |
| Start Scoring | 189.3ms | 5.2 ops/sec | 102 bytes | ⚠️ Good |
| Add Ball | 677.8ms | 1.4 ops/sec | 186 bytes | ❌ Needs Optimization |
| Get Scorecard | 755.2ms | 1.4 ops/sec | 5,749 bytes | ❌ Needs Optimization |
| Complex Scorecard Query | 788.9ms | 1.4 ops/sec | 5,749 bytes | ❌ Needs Optimization |

### Key Findings
- **Series operations are highly optimized** with sub-60ms response times
- **Scorecard operations are the main bottleneck** requiring optimization
- **Data transfer varies significantly** from 240 bytes to 5,749 bytes
- **Database contains 425 series records** with room for optimization

### Optimization Recommendations
1. **Immediate**: Optimize scorecard queries and add database indexes
2. **Medium-term**: Implement caching and database schema optimization  
3. **Long-term**: Consider read replicas and microservices architecture

*Benchmark conducted on Apple M1 Pro, macOS 24.6.0, Go 1.23, Supabase PostgreSQL*

## Project Structure

```
backend/
├── config/          # Configuration management
├── database/        # Database client and operations
├── internal/        # Internal application code
│   ├── handlers/    # HTTP request handlers
│   ├── models/      # Data models and DTOs
│   ├── services/    # Business logic services
│   └── repository/  # Data access layer
├── tests/           # Test files
│   └── performance/ # Performance benchmark tests
├── main.go         # Application entry point
├── go.mod          # Go module file
├── env.example     # Environment variables template
└── README.md       # This file
```
