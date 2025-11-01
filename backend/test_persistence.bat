@echo off
echo Testing user registration and data persistence...

echo.
echo 1. Testing user registration:
curl -X POST http://localhost:8080/api/auth/register ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"testuser2\",\"email\":\"test2@example.com\",\"password\":\"password123\",\"nickname\":\"Test User 2\"}"

echo.
echo.
echo 2. Testing user login:
curl -X POST http://localhost:8080/api/auth/login ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"testuser2\",\"password\":\"password123\"}"

echo.
echo.
echo 3. Testing with existing test user:
curl -X POST http://localhost:8080/api/auth/login ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"testuser\",\"password\":\"password\"}"

echo.
echo.
echo Test completed!