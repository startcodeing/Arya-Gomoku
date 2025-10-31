#!/bin/bash

# Enhanced AI Test Runner
# This script provides easy ways to test the enhanced AI

set -e

echo "🎮 Arya-Gomoku Enhanced AI Test Suite"
echo "===================================="

# Function to print usage
print_usage() {
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  demo        Run AI demo"
    echo "  test        Run unit tests"
    echo "  bench       Run benchmarks"
    echo "  coverage    Run tests with coverage"
    echo "  all         Run all tests and benchmarks"
    echo "  clean       Clean test cache"
    echo "  help        Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 demo     # Run the AI demo"
    echo "  $0 test     # Run unit tests"
    echo "  $0 bench    # Run performance benchmarks"
    echo "  $0 all      # Run everything"
}

# Function to run demo
run_demo() {
    echo "🚀 Running Enhanced AI Demo..."
    echo ""
    cd "$(dirname "$0")"
    go run cmd/ai_demo/main.go
}

# Function to run tests
run_tests() {
    echo "🧪 Running Unit Tests..."
    echo ""
    cd "$(dirname "$0")"
    go test ./internal/service/ -v
}

# Function to run benchmarks
run_benchmarks() {
    echo "⚡ Running Performance Benchmarks..."
    echo ""
    cd "$(dirname "$0")"
    go test ./internal/service/ -bench=. -benchmem
}

# Function to run coverage
run_coverage() {
    echo "📊 Running Tests with Coverage..."
    echo ""
    cd "$(dirname "$0")"

    # Run tests with coverage
    go test ./internal/service/ -coverprofile=coverage.out

    # Show coverage percentage
    coverage=$(go tool cover -func=coverage.out | tail -1)
    echo "Coverage: $coverage"

    # Generate HTML report if available
    if command -v open >/dev/null 2>&1; then
        go tool cover -html=coverage.out -o coverage.html
        echo "📄 Coverage report generated: coverage.html"
        echo "💡 Open coverage.html in your browser to view details"
    fi
}

# Function to run all tests
run_all() {
    echo "🎯 Running Complete Test Suite..."
    echo ""

    echo "--- Demo ---"
    run_demo
    echo ""

    echo "--- Unit Tests ---"
    run_tests
    echo ""

    echo "--- Benchmarks ---"
    run_benchmarks
    echo ""

    echo "--- Coverage ---"
    run_coverage

    echo ""
    echo "✅ All tests completed!"
}

# Function to clean cache
clean_cache() {
    echo "🧹 Cleaning test cache..."
    cd "$(dirname "$0")"
    go clean -testcache
    rm -f coverage.out coverage.html
    echo "✅ Cache cleaned!"
}

# Main script logic
case "${1:-help}" in
    "demo")
        run_demo
        ;;
    "test")
        run_tests
        ;;
    "bench")
        run_benchmarks
        ;;
    "coverage")
        run_coverage
        ;;
    "all")
        run_all
        ;;
    "clean")
        clean_cache
        ;;
    "help"|*)
        print_usage
        ;;
esac