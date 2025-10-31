@echo off
REM Enhanced AI Test Runner for Windows
REM This script provides easy ways to test the enhanced AI

setlocal enabledelayedexpansion

echo 🎮 Arya-Gomoku Enhanced AI Test Suite
echo ====================================

REM Function to print usage
if "%1"=="help" goto :print_usage
if "%1"=="" goto :print_usage

REM Main script logic
if "%1"=="demo" goto :run_demo
if "%1"=="test" goto :run_tests
if "%1"=="bench" goto :run_benchmarks
if "%1"=="coverage" goto :run_coverage
if "%1"=="all" goto :run_all
if "%1"=="clean" goto :clean_cache

goto :print_usage

:run_demo
echo 🚀 Running Enhanced AI Demo...
echo.
cd /d "%~dp0"
go run cmd\ai_demo\main.go
goto :end

:run_tests
echo 🧪 Running Unit Tests...
echo.
cd /d "%~dp0"
go test .\internal\service\ -v
goto :end

:run_benchmarks
echo ⚡ Running Performance Benchmarks...
echo.
cd /d "%~dp0"
go test .\internal\service\ -bench=. -benchmem
goto :end

:run_coverage
echo 📊 Running Tests with Coverage...
echo.
cd /d "%~dp0"

REM Run tests with coverage
go test .\internal\service\ -coverprofile=coverage.out

REM Show coverage percentage
for /f "tokens=*" %%i in ('go tool cover -func=coverage.out ^| findstr total:') do set coverage=%%i
echo Coverage: !coverage!

REM Generate HTML report
go tool cover -html=coverage.out -o coverage.html
echo 📄 Coverage report generated: coverage.html
echo 💡 Open coverage.html in your browser to view details
goto :end

:run_all
echo 🎯 Running Complete Test Suite...
echo.

echo --- Demo ---
call :run_demo
echo.

echo --- Unit Tests ---
call :run_tests
echo.

echo --- Benchmarks ---
call :run_benchmarks
echo.

echo --- Coverage ---
call :run_coverage

echo.
echo ✅ All tests completed!
goto :end

:clean_cache
echo 🧹 Cleaning test cache...
cd /d "%~dp0"
go clean -testcache
if exist coverage.out del coverage.out
if exist coverage.html del coverage.html
echo ✅ Cache cleaned!
goto :end

:print_usage
echo Usage: %0 [command]
echo.
echo Commands:
echo   demo        Run AI demo
echo   test        Run unit tests
echo   bench       Run benchmarks
echo   coverage    Run tests with coverage
echo   all         Run all tests and benchmarks
echo   clean       Clean test cache
echo   help        Show this help message
echo.
echo Examples:
echo   %0 demo     # Run the AI demo
echo   %0 test     # Run unit tests
echo   %0 bench    # Run performance benchmarks
echo   %0 all      # Run everything

:end
pause