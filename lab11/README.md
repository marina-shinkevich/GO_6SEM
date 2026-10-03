# ПИС Лабораторная 11

Состав решения:

- `GO11_01` - REST-сервер на `github.com/gorilla/mux`, функционально повторяющий `GO05_01`.
- `GO11_01\docs` - сгенерированная OpenAPI/Swagger-документация.
- `GO11_01\sql\create_lab11_db.sql` - SQL-скрипт для БД `PIS_Lab5`.
- `GO11_01\requests` - JSON-запросы для проверки.
- `GO11_01\test_go11_01.cmd` - сценарий проверки через `curl.exe`.
- `bin` - сюда собирается исполняемый файл.

## Адреса

- REST API: `http://localhost:3000`
- Swagger UI: `http://localhost:3000/swagger/index.html`

## Подключение к SQL Server

Используется тот же SQL Server, что и в предыдущих лабораторных:

- сервер: `VICTORY\SERVERVILKI`
- порт: `65000`
- база: `PIS_Lab5`
- аутентификация: `Windows Authentication`

## Сборка

```powershell
Set-Location "E:\BSTU\sem6\PIS\lab11\GO11_01"
go build -o ..\bin\GO11_01.exe
```

## Запуск

```powershell
Set-Location "E:\BSTU\sem6\PIS\lab11\bin"
.\GO11_01.exe
```

Проверка в другом окне:

```powershell
Set-Location "E:\BSTU\sem6\PIS\lab11\GO11_01"
.\test_go11_01.cmd
```
