package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func handlerMPF(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "cannot parse multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}


	x := r.FormValue("x")
	y := r.FormValue("y")

	// 3. Извлечение переданного файла (в примере из лекции файл передается под ключом "s").
	// Метод r.FormFile возвращает сам поток файла и объект header с метаданными.
	file, header, err := r.FormFile("s")
	if err != nil {
		http.Error(w, "cannot get file field 's': "+err.Error(), http.StatusBadRequest)
		return
	}
	// Обязательно закрываем дескриптор файла в конце функции, чтобы избежать утечек памяти
	defer file.Close()

	// 4. Чтение содержимого файла (чтение первых 100 байт)
	buf := make([]byte, 100) // создаем небольшой буфер
	n, err := file.Read(buf) // n - количество реально прочитанных байт

	// 5. Формируем ответ клиенту.
	// Выводим текстовые переменные, имя и размер файла из заголовков,
	// а также пытаемся вывести первые 2 байта файла (BOM) в шестнадцатеричном виде (%x)
	// и оставшийся текст файла.
	fmt.Fprintf(w, "x = %s, y = %s, FileName = %s, FileSize = %d, BOM = %x, File content = %s\n",
		x, y, header.Filename, header.Size, string(buf[0:2]), string(buf[2:n]))
}

func main() {
	r := mux.NewRouter()

	// Регистрируем маршрут строго для POST-запроса
	r.HandleFunc("/MPF", handlerMPF).Methods(http.MethodPost)

	log.Println("Server is running on port 3000")
	// Запускаем сервер
	log.Fatal(http.ListenAndServe(":3000", r))
}
