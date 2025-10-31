# Enhanced AI Testing Guide

This guide explains how to test the enhanced AI with minimax + alpha-beta pruning implementation.

## Quick Start

### Option 1: Use the Test Scripts (Recommended)

**Windows:**
```cmd
# Run demo
test_ai.bat demo

# Run all tests
test_ai.bat all

# See all options
test_ai.bat help
```

**Linux/macOS:**
```bash
# Make script executable
chmod +x test_ai.sh

# Run demo
./test_ai.sh demo

# Run all tests
./test_ai.sh all

# See all options
./test_ai.sh help
```

### Option 2: Manual Commands

```bash
# Run the AI demo
go run cmd/ai_demo/main.go

# Run unit tests
go test ./internal/service/ -v

# Run benchmarks
go test ./internal/service/ -bench=. -benchmem

# Run tests with coverage
go test ./internal/service/ -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Test Components

### 1. AI Demo (`cmd/ai_demo/main.go`)
An interactive demonstration of AI capabilities:
- Empty board analysis
- Win detection scenarios
- Block detection scenarios
- Difficulty level comparison
- Performance statistics

### 2. Unit Tests (`internal/service/enhanced_ai_service_test.go`)
Comprehensive unit tests covering:
- Move generation on empty boards
- Win/block detection
- Time limit enforcement
- Board evaluation algorithms
- Transposition table functionality
- Difficulty level validation
- Statistics tracking

### 3. Benchmark Tests (`internal/service/enhanced_ai_benchmark_test.go`)
Performance benchmarks for:
- AI move generation speed
- Complex board analysis
- Individual algorithm components
- Memory usage patterns
- Stress testing scenarios

## Running Specific Tests

### Demo Scenarios
```bash
go run cmd/ai_demo/main.go
```

### Unit Test Categories
```bash
# All tests
go test ./internal/service/ -v

# Specific tests
go test ./internal/service/ -run TestGetAIMove_EmptyBoard -v
go test ./internal/service/ -run TestCheckWin -v
go test ./internal/service/ -run TestTranspositionTable -v
```

### Benchmark Categories
```bash
# All benchmarks
go test ./internal/service/ -bench=. -benchmem

# Specific benchmarks
go test ./internal/service/ -bench=BenchmarkGetAIMove -benchmem
go test ./internal/service/ -bench=BenchmarkTranspositionTable -benchmem
```

## Understanding Test Results

### Demo Output
The demo shows:
- Move coordinates chosen by AI
- Evaluation scores
- Time taken per move
- Performance statistics

### Unit Test Results
- ✅ Pass: All assertions passed
- ❌ Fail: One or more assertions failed
- 📊 Coverage: Percentage of code tested

### Benchmark Results
```
BenchmarkGetAIMove/Medium-8            	 2000	    650000 ns/op	 512 B/op	       8 allocs/op
```
- `2000`: Number of iterations run
- `650000 ns/op`: Nanoseconds per operation (0.65ms)
- `512 B/op`: Bytes allocated per operation
- `8 allocs/op`: Memory allocations per operation

## Expected Performance Characteristics

Based on testing on modern hardware:

| Difficulty | Response Time | Memory Usage | Pruning Efficiency |
|------------|---------------|--------------|-------------------|
| Easy       | < 100ms       | ~1MB         | 40-50%            |
| Medium     | 100-500ms     | ~2MB         | 50-60%            |
| Hard       | 500ms-2s      | ~5MB         | 55-65%            |
| Expert     | 2-5s          | ~10MB        | 60-70%            |

## Troubleshooting

### Common Issues

**Build Errors:**
```bash
# Ensure you're in the backend directory
cd backend

# Check Go version (requires Go 1.21+)
go version

# Install dependencies
go mod tidy
```

**Test Timeouts:**
- Expert difficulty may take longer on slower machines
- Use `test_ai.bat bench` instead of `test_ai.bat all` for faster testing

**Memory Issues:**
- Clear transposition table: `test_ai.bat clean`
- Reduce test iterations in benchmark code

### Performance Tips

1. **First Run:** First-time runs may be slower due to compilation
2. **Cache Warming:** Subsequent moves are faster due to transposition table
3. **System Load:** Close other applications for consistent benchmarking
4. **Power Settings:** Use high-performance power settings

## Integration with Application

To integrate tests into your development workflow:

### Pre-commit Testing
```bash
# Quick validation
go test ./internal/service/ -run TestGetAIMove_EmptyBoard
```

### CI/CD Pipeline
```yaml
# Example GitHub Actions step
- name: Test Enhanced AI
  run: |
    cd backend
    go test ./internal/service/ -v
    go test ./internal/service/ -bench=BenchmarkGetAIMove -benchmem
```

### Performance Monitoring
```bash
# Regular performance checks
go test ./internal/service/ -bench=BenchmarkGetAIMove/Medium -benchmem
```

## Advanced Testing

### Custom Scenarios
You can modify the demo or add custom tests by:
1. Editing `cmd/ai_demo/main.go` for new scenarios
2. Adding test functions to `enhanced_ai_service_test.go`
3. Creating custom benchmarks in `enhanced_ai_benchmark_test.go`

### Stress Testing
```bash
# Run extended stress test
go test ./internal/service/ -run TestStressTest -v -count=5
```

### Memory Profiling
```bash
# Generate memory profile
go test ./internal/service/ -bench=BenchmarkGetAIMove -memprofile=mem.prof
go tool pprof mem.prof
```

## Contributing

When adding new AI features:
1. Add corresponding unit tests
2. Include performance benchmarks
3. Update demo if relevant
4. Test all difficulty levels
5. Verify backward compatibility

## Support

If you encounter issues:
1. Check Go version compatibility
2. Clean test cache: `test_ai.bat clean`
3. Run individual tests to isolate problems
4. Check system resources (memory, CPU)

For detailed implementation information, see `ENHANCED_AI_GUIDE.md`.