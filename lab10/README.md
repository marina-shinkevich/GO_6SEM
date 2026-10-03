# ПИС Лабораторная 10

Состав решения:

- `GO10_01` - web-сервер на `GraphQL` (`graph-gophers/graphql-go`), функционально повторяющий `GO05_01` и `GO06_01`.
- `GO10_01\sql\create_lab10_db.sql` - SQL-скрипт для общей БД `PIS_Lab5`, таблицы и стартовых данных.
- `GO10_01\requests` - JSON-запросы GraphQL для проверки.
- `GO10_01\test_go10_01.cmd` - сценарий проверки через `curl.exe`.
- `bin` - сюда собирается исполняемый файл.

## GraphQL endpoint

Сервер слушает порт `3000`.

- GraphQL: `http://localhost:3000/graphql`
- GET-заглушка: `http://localhost:3000/`

Поддерживаемые операции:

- `celebrities` - вся коллекция.
- `celebrity(id: Int!)` - элемент по `id`.
- `addCelebrity(input: AddCelebrityInput!)` - добавление элемента.
- `updateCelebrity(id: Int!, input: UpdateCelebrityInput!)` - изменение элемента.
- `deleteCelebrity(id: Int!)` - удаление элемента.

## Подключение к SQL Server

Используется та же база, что в лабораторных 5 и 6:

- сервер в SSMS: `VICTORY\SERVERVILKI`
- сервер в Go-приложении: `tcp:VICTORY`, порт `65000`
- база: `PIS_Lab5`
- аутентификация: `Windows Authentication`
- шифрование: `Encrypt=True`
- доверять сертификату сервера: `True`

## Создание БД вручную

```powershell
Set-Location "E:\BSTU\sem6\PIS\lab10\GO10_01"
sqlcmd -S "VICTORY\SERVERVILKI" -E -C -i .\sql\create_lab10_db.sql
```

Приложение также создаёт базу и таблицу автоматически при старте, если их ещё нет.

## Сборка

```powershell
Set-Location "E:\BSTU\sem6\PIS\lab10\GO10_01"
go build -o ..\bin\GO10_01.exe
```

## Запуск

По умолчанию приложение подключается к:

- `GO10_SQLSERVER=VICTORY`
- `GO10_SQLPORT=65000`
- `GO10_DATABASE=PIS_Lab5`

Если нужно, параметры можно переопределить переменными окружения:

```powershell
$env:GO10_SQLSERVER = "VICTORY"
$env:GO10_SQLPORT = "65000"
$env:GO10_DATABASE = "PIS_Lab5"
```

Запуск:

```powershell
Set-Location "E:\BSTU\sem6\PIS\lab10\bin"
.\GO10_01.exe
```

Потом в другом окне:

```powershell
Set-Location "E:\BSTU\sem6\PIS\lab10\GO10_01"
.\test_go10_01.cmd
```
