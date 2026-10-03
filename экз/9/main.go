package main

import (
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/gorilla/mux"
)

// handlerDFS (Download File Server) — функция для скачивания файла.
// Ожидает запрос вида: GET http://localhost:3000/DFS/static/PNorton.jpg
func handlerDFS(w http.ResponseWriter, r *http.Request) {
	// 1. Указываем базовую директорию на сервере, где лежат файлы
	const baseDir = "./static"

	// 2. Извлекаем имя запрошенного файла из URL-шаблона
	vars := mux.Vars(r)
	filename := vars["filename"]

	// 3. Безопасно склеиваем пути (базовую директорию и имя файла)
	// filepath.Join автоматически ставит правильные слеши для Windows или Linux
	filePath := filepath.Join(baseDir, filename)

	// 4. Получаем информацию о файле с помощью os.Stat
	info, err := os.Stat(filePath)

	if err == nil {
		// Если файл найден, проверяем, не является ли он директорией
		if !info.IsDir() {
			// Получаем чистое имя файла (без путей) для отправки клиенту
			downloadName := path.Base(filename)

			// === КЛЮЧЕВОЙ БЛОК ДЛЯ СКАЧИВАНИЯ ===
			// Заголовок Content-Disposition со значением "attachment"
			// заставляет браузер именно СКАЧАТЬ файл, а не пытаться его открыть.
			w.Header().Set("Content-Disposition", "attachment; filename="+downloadName)

			// Заголовок application/octet-stream говорит, что это произвольный бинарный файл
			w.Header().Set("Content-Type", "application/octet-stream")
			// =====================================

			// 5. Встроенная функция Go, которая берет файл с диска и безопасно
			// перекачивает его байты в сетевой поток ответа клиенту
			http.ServeFile(w, r, filePath)
		} else {
			// Если по этому пути лежит папка, а не файл
			http.Error(w, "not a file", http.StatusBadRequest)
		}
	} else if os.IsNotExist(err) {
		// Если файл физически не существует на диске
		http.Error(w, "file not found", http.StatusNotFound)
	} else {
		// Если произошла другая системная ошибка (например, нет прав доступа)
		http.Error(w, "stat error: "+err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	r := mux.NewRouter()

	// Регистрируем маршрут с шаблоном {filename}
	// Любой текст после /DFS/static/ будет сохранен в переменную filename
	r.HandleFunc("/DFS/static/{filename}", handlerDFS).Methods(http.MethodGet)

	log.Println("Server is running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
