@echo off
echo Testing GET /A
curl -X GET http://localhost:3000/A
echo.
echo Testing POST /A/B
curl -X POST http://localhost:3000/A/B
echo.
echo Testing PUT /Unknown
curl -X PUT http://localhost:3000/Unknown
pause