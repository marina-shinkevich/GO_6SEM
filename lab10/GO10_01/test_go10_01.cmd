@echo off
setlocal


echo.
echo QUERY celebrities
curl.exe -s -i -X POST http://localhost:3000/graphql ^
  -H "Content-Type: application/json" ^
  --data-binary "@requests/all.json"
echo.

echo.
echo QUERY celebrity id=2
curl.exe -s -i -X POST http://localhost:3000/graphql ^
  -H "Content-Type: application/json" ^
  --data-binary "@requests/one.json"
echo.

echo.
echo MUTATION addCelebrity
curl.exe -s -i -X POST http://localhost:3000/graphql ^
  -H "Content-Type: application/json" ^
  --data-binary "@requests/create.json"
echo.

echo.
echo MUTATION addCelebrity duplicate
curl.exe -s -i -X POST http://localhost:3000/graphql ^
  -H "Content-Type: application/json" ^
  --data-binary "@requests/create.json"
echo.

echo.
echo MUTATION updateCelebrity id=2
curl.exe -s -i -X POST http://localhost:3000/graphql ^
  -H "Content-Type: application/json" ^
  --data-binary "@requests/update.json"
echo.

echo.
echo MUTATION deleteCelebrity id=3
curl.exe -s -i -X POST http://localhost:3000/graphql ^
  -H "Content-Type: application/json" ^
  --data-binary "@requests/delete.json"
echo.

echo.
echo MUTATION deleteCelebrity id=30
curl.exe -s -i -X POST http://localhost:3000/graphql ^
  -H "Content-Type: application/json" ^
  --data-binary "@requests/delete_missing.json"
echo.

endlocal
Pause
