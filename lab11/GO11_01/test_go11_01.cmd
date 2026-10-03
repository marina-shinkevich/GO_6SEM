@echo off
setlocal


echo.
echo GET /Celebrities/All
curl.exe -s -i http://localhost:3000/Celebrities/All
echo.

echo.
echo GET /Celebrities/2
curl.exe -s -i http://localhost:3000/Celebrities/2
echo.

echo.
echo POST /Celebrities
curl.exe -s -i -X POST http://localhost:3000/Celebrities ^
  -H "Content-Type: application/json" ^
  --data-binary "@requests/create.json"
echo.

echo.
echo POST /Celebrities duplicate
curl.exe -s -i -X POST http://localhost:3000/Celebrities ^
  -H "Content-Type: application/json" ^
  --data-binary "@requests/create.json"
echo.

echo.
echo PUT /Celebrities/2
curl.exe -s -i -X PUT http://localhost:3000/Celebrities/2 ^
  -H "Content-Type: application/json" ^
  --data-binary "@requests/update.json"
echo.

echo.
echo DELETE /Celebrities/3
curl.exe -s -i -X DELETE http://localhost:3000/Celebrities/3
echo.

echo.
echo DELETE /Celebrities/30
curl.exe -s -i -X DELETE http://localhost:3000/Celebrities/30
echo.

endlocal
Pause
