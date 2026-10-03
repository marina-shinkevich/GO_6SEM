package main

import (
	"strconv"

	// Подключаем фреймворк Iris
	"github.com/kataras/iris/v12"

	// Подключаем обертки Swagger специально для фреймворка Iris
	"github.com/iris-contrib/swagger/v12"
	"github.com/iris-contrib/swagger/v12/swaggerFiles"

	// ВАЖНО: Подключаем сгенерированную папку docs!
	_ "rest/docs"
)

// === 1. МОДЕЛИ ДАННЫХ И ХРАНИЛИЩЕ ===

// Celebrity - структура данных
// @Description Модель знаменитости.
type Celebrity struct {
	// @Description Уникальный идентификатор знаменитости
	Id int `json:"id"`

	// @Description Полное Имя знаменитости (обязательное поле)
	FullName string `json:"fullname"`
}

var (
	// Наше локальное in-memory хранилище
	celebrities = []Celebrity{
		{Id: 1, FullName: "John Doe"},
		{Id: 2, FullName: "Jane Smith"},
	}
	// Счетчик для автоинкремента ID при создании новых записей
	nextID = 3
)

// === 2. ОБЩЕЕ ОПИСАНИЕ API ===

// @title Celebrities API
// @version 1.0
// @description Полное REST API для управления списком знаменитостей с использованием локального хранилища.
// @host localhost:3001
// @BasePath /
func main() {
	app := iris.New()

	// Swagger UI
	app.Get("/docs/{any:path}", swagger.WrapHandler(swaggerFiles.Handler))

	// Группируем маршруты REST API
	cRoute := app.Party("/celebrities")
	{
		cRoute.Get("/", getAllCelebrities)          // Получить всех
		cRoute.Get("/{id:int}", getCelebrityByID)   // Получить одного по ID
		cRoute.Post("/", createCelebrity)           // Создать нового
		cRoute.Put("/{id:int}", updateCelebrity)    // Обновить существующего
		cRoute.Delete("/{id:int}", deleteCelebrity) // Удалить по ID
	}

	app.Listen(":3001")
}

// === 3. ОБРАБОТЧИКИ REST ОПЕРАЦИЙ ===

// getAllCelebrities godoc
// @Summary Получить всех знаменитостей
// @Description Возвращает полный список знаменитостей
// @Produce json
// @Success 200 {array} Celebrity
// @Router /celebrities [get]
func getAllCelebrities(ctx iris.Context) {
	ctx.JSON(celebrities)
}

// getCelebrityByID godoc
// @Summary Получить знаменитость по ID
// @Description Возвращает одну знаменитость по её числовому идентификатору
// @Produce json
// @Param id path int true "ID Знаменитости"
// @Success 200 {object} Celebrity
// @Failure 404 {string} string "Celebrity not found"
// @Router /celebrities/{id} [get]
func getCelebrityByID(ctx iris.Context) {
	id, _ := ctx.Params().GetInt("id")

	for _, c := range celebrities {
		if c.Id == id {
			ctx.JSON(c)
			return
		}
	}

	ctx.StatusCode(iris.StatusNotFound)
	ctx.WriteString("Celebrity not found")
}

// createCelebrity godoc
// @Summary Создать новую знаменитость
// @Description Принимает JSON без ID, генерирует ID автоматически и сохраняет в массив
// @Accept json
// @Produce json
// @Param celebrity body Celebrity true "Данные знаменитости"
// @Success 201 {object} Celebrity
// @Failure 400 {string} string "Bad Request"
// @Router /celebrities [post]
func createCelebrity(ctx iris.Context) {
	var newCelebrity Celebrity
	if err := ctx.ReadJSON(&newCelebrity); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("Invalid JSON input")
		return
	}

	newCelebrity.Id = nextID
	nextID++
	celebrities = append(celebrities, newCelebrity)

	ctx.StatusCode(iris.StatusCreated)
	ctx.JSON(newCelebrity)
}

// updateCelebrity godoc
// @Summary Обновить существующую знаменитость
// @Description Находит знаменитость по ID и перезаписывает её данные
// @Accept json
// @Produce json
// @Param id path int true "ID Знаменитости"
// @Param celebrity body Celebrity true "Новые данные знаменитости"
// @Success 200 {object} Celebrity
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Celebrity not found"
// @Router /celebrities/{id} [put]
func updateCelebrity(ctx iris.Context) {
	id, _ := ctx.Params().GetInt("id")
	var updatedCelebrity Celebrity

	if err := ctx.ReadJSON(&updatedCelebrity); err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("Invalid JSON input")
		return
	}

	for i, c := range celebrities {
		if c.Id == id {
			celebrities[i].FullName = updatedCelebrity.FullName
			updatedCelebrity.Id = id // Возвращаем клиенту объект с верным ID
			ctx.JSON(celebrities[i])
			return
		}
	}

	ctx.StatusCode(iris.StatusNotFound)
	ctx.WriteString("Celebrity not found")
}

// deleteCelebrity godoc
// @Summary Удалить знаменитость
// @Description Удаляет знаменитость из локального массива по ID
// @Produce json
// @Param id path int true "ID Знаменитости"
// @Success 200 {string} string "Celebrity deleted successfully"
// @Failure 404 {string} string "Celebrity not found"
// @Router /celebrities/{id} [delete]
func deleteCelebrity(ctx iris.Context) {
	id, _ := ctx.Params().GetInt("id")

	for i, c := range celebrities {
		if c.Id == id {
			// Удаляем элемент из слайса
			celebrities = append(celebrities[:i], celebrities[i+1:]...)
			ctx.WriteString("Celebrity with ID " + strconv.Itoa(id) + " deleted successfully")
			return
		}
	}

	ctx.StatusCode(iris.StatusNotFound)
	ctx.WriteString("Celebrity not found")
}