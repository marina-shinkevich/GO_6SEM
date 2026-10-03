@echo off
echo Сборка модуля GO06_01 с использованием GORM...
echo.

echo 1. Очистка предыдущей сборки...
if exist bin\GO06_01.exe del bin\GO06_01.exe

echo 2. Загрузка зависимостей...
go mod tidy

echo 3. Сборка приложения...
go build -o bin\GO06_01.exe main_gorm.go

echo 4. Проверка сборки...
if exist bin\GO06_01.exe (
    echo Сборка успешно завершена!
    echo Исполняемый файл: bin\GO06_01.exe
    echo.
    echo Для запуска сервера выполните:
    echo   bin\GO06_01.exe
) else (
    echo Ошибка при сборке!
    exit /b 1
)