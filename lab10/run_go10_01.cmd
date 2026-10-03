@echo off
setlocal

for /f "tokens=5" %%P in ('netstat -ano -p tcp ^| findstr /R /C:":3000 .*LISTENING"') do (
    echo Port 3000 is already busy by PID %%P.
    echo Stop it with:
    echo   taskkill /PID %%P /F
    echo.
    pause
    exit /b 1
)

cd /d "%~dp0bin"
GO10_01.exe

echo.
echo Program finished with exit code %ERRORLEVEL%.
pause
