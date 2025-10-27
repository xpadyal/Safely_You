# Safely You - Device Monitoring API

A high-performance Go-based REST API for monitoring device connectivity and performance metrics. This service provides real-time device uptime tracking and upload performance analytics for IoT and safety-critical applications.

## Features

- **Real-time Device Monitoring**: Track device connectivity with precise uptime calculations
- **Performance Analytics**: Monitor upload times and network performance metrics
- **Thread-Safe Architecture**: Concurrent request handling with proper synchronization
- **Robust Input Validation**: Timestamp validation with reasonableness checks
- **CSV Device Management**: Bulk device initialization from configuration files
- **Simulator Integration**: Seamless testing with provided device simulator

## API Endpoints

### Health Check
```
GET /health
```
Returns server health status.

### Device Heartbeat
```
POST /api/v1/devices/{device_id}/heartbeat
Content-Type: application/json

{
  "sent_at": "2023-10-25T19:47:31Z"
}
```

### Upload Statistics
```
POST /api/v1/devices/{device_id}/stats
Content-Type: application/json

{
  "sent_at": "2023-10-25T19:47:31Z",
  "upload_time": 5000000000
}
```

### Get Device Statistics
```
GET /api/v1/devices/{device_id}/stats
```

Response:
```json
{
  "uptime": 85.67,
  "avg_upload_time": "5.2s"
}
```

## Prerequisites

- **Go 1.19+** (tested with Go 1.21)
- **Device Simulator**: `device-sim` executable (included)
- **Device Configuration**: `devices.csv` file with device definitions

## Quick Start

1. **Clone and setup**
   ```bash
   git clone https://github.com/xpadyal/Safely_You.git
   cd Safely_You
   go mod download
   chmod +x device-sim
   ```

2. **Verify configuration**
   ```bash
   ls -la devices.csv device-sim  # Ensure both files exist
   ```

## Usage

### Development Mode

**Start the server:**
```bash
go run ./cmd/server
```

**Test the API:**
```bash
# Health check
curl http://localhost:8080/health

# Send heartbeat
curl -X POST http://localhost:8080/api/v1/devices/{device-id}/heartbeat \
  -H "Content-Type: application/json" \
  -d '{"sent_at": "2023-10-25T19:47:31Z"}'

# Send upload stats
curl -X POST http://localhost:8080/api/v1/devices/{device-id}/stats \
  -H "Content-Type: application/json" \
  -d '{"sent_at": "2023-10-25T19:47:31Z", "upload_time": 5000000000}'

# Get device statistics
curl http://localhost:8080/api/v1/devices/{device-id}/stats
```

### With Device Simulator

**Start the server:**
```bash
go run ./cmd/server
# Output: Loaded devices from devices.csv
#         [GIN-debug] Listening and serving HTTP on :8080
```

**Run the simulator:**
```bash
./device-sim -port 8080
# Simulator will:
# - Send heartbeat and upload data to all devices
# - Query statistics for each device
# - Output results to results.txt and console
```

**View results:**
```bash
cat results.txt
```

### Production Build

```bash
# Build optimized binary
go build -o Safely_You ./cmd/server

# Run server
./Safely_You

# Run simulator
./device-sim
```

## Configuration

Configure the server using environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `:8080` | Server port |
| `GIN_MODE` | `debug` | Gin framework mode |

**Example:**
```bash
export PORT=:9000
export GIN_MODE=release
go run ./cmd/server
```

## Project Structure

```
Safely_You/
├── cmd/server/           # Application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── handlers/        # HTTP request handlers
│   ├── loader/          # CSV data loading
│   ├── models/          # Data structures
│   ├── store/           # Business logic and data operations
│   ├── utils/           # Utility functions
│   └── validation/      # Input validation helpers
├── tests/               # Unit tests
├── devices.csv          # Device definitions
├── device-sim           # Device simulator executable
└── go.mod              # Go module definition
```

## Uptime Calculation

The uptime calculation uses an **inclusive span** methodology:

- **Formula**: `(unique_minute_heartbeats / total_minute_span) × 100`
- **Span**: Includes all minute buckets from first to last heartbeat
- **Example**: Heartbeats at 10:00 and 10:02 span 3 minutes (10:00, 10:01, 10:02)
- **Result**: 2 heartbeats ÷ 3 minutes = 66.67% uptime

This approach ensures missed heartbeats negatively impact uptime, providing accurate reliability metrics for safety-critical applications.

## API Response Codes

| Code | Description |
|------|-------------|
| `200 OK` | Successful request with response data |
| `204 No Content` | Successful operation, no response body |
| `400 Bad Request` | Invalid input data or malformed request |
| `404 Not Found` | Device not found or invalid device ID |
| `500 Internal Server Error` | Server-side error or processing failure |

## Thread Safety

The application implements thread-safe operations using `sync.RWMutex`:

- **Read Operations**: Concurrent access with read locks for optimal performance
- **Write Operations**: Exclusive access with write locks for data integrity
- **Device Operations**: All device data access is properly synchronized
- **Concurrent Requests**: Supports multiple simultaneous API calls safely

## Testing

**Run the test suite:**
```bash
go test ./tests/...
```

**Run with coverage:**
```bash
go test -cover ./tests/...
```

## Performance Characteristics

- **Response Times**: 1-100μs for heartbeat operations, 50-100μs for stats queries
- **Concurrent Access**: Supports multiple simultaneous requests with read/write locks
- **Memory Efficiency**: Uses slices and maps for optimal performance
- **Input Validation**: Timestamp reasonableness checks (24h past, 5min future)
- **Scalability**: Handles 10,000+ devices and 100,000+ heartbeats per device

## Troubleshooting

### Common Issues

**Port already in use:**
```bash
# Kill existing process
pkill -f Safely_You
# Or use different port
export PORT=:9000
```

**Device simulator 404 errors:**
- Ensure server is running on port 8080
- Verify API routes include `/api/v1` prefix

**CSV loading errors:**
- Verify `devices.csv` exists in project root
- Check CSV format (header row + device IDs)

### Request Logging

The server provides detailed request logging:
```
[GIN] 2025/10/25 - 19:47:31 | 204 | 2.208µs | 127.0.0.1 | POST "/api/v1/devices/b4-45-52-a2-f1-3c/heartbeat"
```

Format: `[timestamp] | status | duration | client | method path`

## Technical Write-Up

### 1) Time Spent & Most Difficult Part

**Total Time**: About 4-6 hours.

**Breakdown**:
- **Learning Go (1.5h)**: Coming from Python, I had to wrap my head around Go's syntax and paradigms. I actually fell in love with Go's explicit concurrency model - it's so much clearer than Python's threading. The static typing felt restrictive at first, but I quickly appreciated how it catches errors at compile time rather than runtime.
- **Understanding Requirements (30min)**: The OpenAPI spec was pretty clear, but the inclusive vs exclusive span calculation for uptime took some thinking. I initially implemented it wrong and had to rework it.
- **Architecture Planning (45min)**: Deciding on package structure and how to handle thread safety. I wanted to keep it simple but well-organized.
- **Core Implementation (2h)**: Building the store operations, uptime calculation, and API handlers. Getting Gin set up was straightforward.
- **Testing & Integration (1h)**: Getting the device simulator working and making sure everything matched the OpenAPI spec exactly.
- **Refactoring (45min)**: Adding proper error handling, implementing thread safety, and organizing everything into proper packages.

**Most Difficult Part**: Definitely the uptime calculation. I initially calculated it as an exclusive span (only counting minutes where heartbeats actually occurred), but I realized that missed heartbeats should hurt the uptime percentage. Switching to an inclusive span calculation required careful thinking about edge cases and making sure the math was correct.

**Extra Improvement**: After implementing the minimal working version, I added thread safety with `sync.RWMutex` to handle concurrent requests properly. This was tricky - I had to make sure I was using read/write locks correctly without creating deadlocks or hurting performance.

### 2) Development Approach

I took a **bottom-up approach** - starting with the core data structures and business logic, then building the API layer on top. This felt natural coming from Python where I'm used to thinking about data first.

**Key Decisions**:
- **In-memory storage**: Chose simplicity over scalability for this assessment
- **Package structure**: Separated concerns into logical packages (`handlers`, `models`, `store`, etc.)
- **Thread safety**: Added `sync.RWMutex` early to avoid concurrency issues
- **Validation layer**: Centralized all input validation in one place
- **Error handling**: Used helper functions to avoid repetitive error response code

I prioritized **correctness over performance** - making sure the uptime calculation was mathematically sound and the API matched the OpenAPI spec.

### 3) Extensibility — Adding New Metrics

The current architecture makes adding new metrics pretty straightforward:

**Data Model**: Just extend the `Device` struct with new fields:
```go
type Device struct {
    Heartbeats     []time.Time
    UploadTimes    []int64
    CpuUsage       []float64    // New metric
    MemoryUsage    []int64      // New metric
    NetworkLatency []int64      // New metric
}
```

**Handler Pattern**: Create metric-specific handlers that follow the same validation and storage pattern as the existing upload stats handler.

**Storage Strategy**: For small scale, the current approach works fine. For larger deployments, I'd migrate to a time-series database like InfluxDB or TimescaleDB.

**API Versioning**: I'd implement proper versioning strategy to maintain backward compatibility as the API evolves.

### 4) Runtime & Space Complexity

**Time Complexity**:
- **Data Insertion**: O(1) - Just appending to slices
- **Device Lookup**: O(1) - Hash map access
- **Uptime Calculation**: O(n) - Linear scan through heartbeats to find unique minute buckets
- **Average Upload Time**: O(m) - Linear scan through upload times

**Space Complexity**:
- **Storage**: O(n + m) - Linear growth with total heartbeats and upload times
- **Temporary Calculations**: O(k) - Proportional to unique minute buckets during uptime calculation

The solution scales reasonably well for moderate datasets. I prioritized correctness over micro-optimizations.

### 5) Performance Characteristics & Bottlenecks

**Current Performance**:
- Response times: 1-100μs for heartbeat operations, 50-100μs for stats queries
- Supports concurrent requests with read/write locks
- Estimated capacity: 10,000+ devices, 100,000+ heartbeats per device

**Identified Bottlenecks**:
- **Recalculating uptime on every request**: Could be expensive with lots of data - caching would help
- **Unbounded memory growth**: No data retention policy - memory keeps growing
- **Single-threaded calculations**: Could be parallelized for very large datasets

### 6) Team Collaboration & Code Organization

**What I implemented**:
- **Modular architecture**: Clean separation into logical packages
- **Handler organization**: Separate files for each endpoint (`heartbeat.go`, `stats.go`, `health.go`)
- **Centralized validation**: All validation logic in one place for consistency
- **Configuration management**: Environment-based config for flexible deployment
- **Clear interfaces**: Well-defined store operations that are easy to test and mock

**For production teams, I'd add**:
- Comprehensive testing (unit, integration, contract tests)
- API versioning strategy
- OpenAPI documentation
- Structured logging and metrics
- CI/CD pipeline with automated testing
- Code review guidelines and quality gates

### 7) Production Readiness Assessment

I focused on getting the core functionality right, but I know there are some **major production safety concerns** I didn't address:

**Security**: No authentication, rate limiting, or input sanitization. In production, this would need JWT auth, rate limiting middleware, and proper input validation to prevent attacks.

**Memory Management**: The biggest concern is unbounded memory growth. With no data retention policy, the application will eventually run out of memory. Production would need TTL-based cleanup or migration to a time-series database.

**Error Handling**: Basic error responses without detailed logging. Production would need structured logging, error tracking (like Sentry), and proper monitoring to debug issues.

**Monitoring**: Only basic Gin logging. Production would need Prometheus metrics, health checks, and alerting to detect problems before they become critical.


**Assessment Philosophy**: I focused on demonstrating solid engineering fundamentals - clean code, proper architecture, and correct business logic - rather than trying to build a production-ready system in a few hours.

### 9) License / Assessment Context

This project is part of a technical assessment for Safely You, demonstrating my ability to build a device monitoring API with Go.

---

**Thank you for taking the time to review this project!** I hope this demonstrates my approach to problem-solving, technical decision-making, and ability to learn new technologies quickly. I'm excited about the opportunity to contribute to Safely You's mission.

