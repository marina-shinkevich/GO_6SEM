# Модуль GO06_01 - Веб-сервер с использованием GORM (Исправленная версия)

## Исправленные проблемы

### 1. **Конфликт имен таблиц и колонок**
**Проблема**: GORM пытался создать таблицу `celebrities` (lowercase), в то время как в базе уже существовала таблица `Celebrities` (uppercase).

**Решение**: 
- Добавлен метод `TableName()` для явного указания имени таблицы
- Указаны точные имена колонок через тег `gorm:"column:ИмяКолонки"`

```go
func (Celebrity) TableName() string {
    return "Celebrities"
}

type Celebrity struct {
    ID           int    `json:"id" gorm:"primaryKey;column:Id"`
    // ...
}
```

### 2. **Дублирование колонки ID**
**Проблема**: GORM пытался добавить колонку `id` как `bigint IDENTITY(1,1)`, но она уже существовала как `INT PRIMARY KEY`.

**Решение**: Заменен `AutoMigrate()` на ручную проверку и создание таблицы:

```go
if !db.Migrator().HasTable(&Celebrity{}) {
    // Создание таблицы вручную
    createTableSQL := `...`
    db.Exec(createTableSQL)
}
```

### 3. **Медленные SQL запросы**
**Проблема**: GORM выполнял медленные запросы для анализа структуры таблицы (276ms и 674ms).

**Решение**: Настроен кастомный логгер GORM:
- Увеличен `SlowThreshold` до 1 секунды
- Установлен `LogLevel` на `Warn` (только предупреждения и ошибки)
- Включен `IgnoreRecordNotFoundError`

### 4. **Конфликт автоинкремента**
**Проблема**: При использовании `db.Create()` GORM пытался использовать автоинкремент для поля ID.

**Решение**: Для операции INSERT используется явный SQL запрос:
```go
db.Exec("INSERT INTO Celebrities (Id, FullName, ...) VALUES (?, ?, ...)", ...)
```

## Обновленная конфигурация GORM

```go
newLogger := logger.New(
    log.New(log.Writer(), "\r\n", log.LstdFlags),
    logger.Config{
        SlowThreshold:             time.Second,
        LogLevel:                  logger.Warn,
        IgnoreRecordNotFoundError: true,
        Colorful:                  false,
    },
)

db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{
    Logger: newLogger,
})
```

## Запуск исправленной версии

1. **Очистка предыдущей сборки**:
   ```bash
   del bin\GO06_01.exe
   ```

2. **Сборка**:
   ```bash
   go build -o bin\GO06_01.exe main_gorm.go
   ```

3. **Запуск**:
   ```bash
   bin\GO06_01.exe
   ```

## Ожидаемое поведение

1. **При первом запуске** (таблица не существует):
   ```
   Подключение к базе данных установлено
   Таблица Celebrities не существует, создание...
   Таблица Celebrities создана
   Сервер запущен на порту 3000
   ```

2. **При последующих запусках** (таблица существует):
   ```
   Подключение к базе данных установлено
   Таблица Celebrities уже существует
   Сервер запущен на порту 3000
   ```

## Тестирование

Все эндпоинты остались прежними:
- `GET /Celebrities/All`
- `GET /Celebrities/{id}`
- `POST /Celebrities`
- `PUT /Celebrities/{id}`
- `DELETE /Celebrities/{id}`

Функциональность полностью идентична модулю GO05_01, но с использованием GORM для работы с базой данных.