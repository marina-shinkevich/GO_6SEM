package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	baseURL    = "http://localhost:3000/"
	destURL    = "http://localhost:3000/"
	baseFolder = "A"
	workFolder = "B"
)

// request — универсальная функция для отправки любых HTTP/WebDAV запросов
func request(method, url string, body io.Reader, h map[string]string) {
	// Создаем запрос с нестандартным методом (например, "MKCOL" или "PROPFIND")
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		log.Println(method, "error:", err)
		return
	}

	// Устанавливаем заголовки (Headers) из переданного словаря
	for k, v := range h {
		req.Header.Set(k, v)
	}

	// Выполняем запрос
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(method, "error:", err)
		return
	}
	defer resp.Body.Close()

	// Читаем и выводим ответ от сервера
	b, _ := io.ReadAll(resp.Body)
	fmt.Println("===", method, url, "===")
	fmt.Println("status:", resp.Status)
	fmt.Println(string(b))
	fmt.Println()
}

func main() {
	// 1. MKCOL (Make Collection) - Создание папки "A"
	request("MKCOL", baseURL+baseFolder+"/", nil, nil)

	// 2. PUT - Загрузка файла hello.txt внутрь папки "A"
	request("PUT", baseURL+baseFolder+"/hello.txt",
		bytes.NewBuffer([]byte("Hello WebDAV")), // Содержимое файла
		map[string]string{"Content-Type": "text/plain"},
	)

	// 3. GET - Скачивание файла (проверка, что он записался)
	request("GET", baseURL+baseFolder+"/hello.txt", nil, nil)

	// 4. PROPFIND - Получение свойств папки "A" в формате XML
	request("PROPFIND", baseURL+baseFolder+"/",
		bytes.NewBuffer([]byte(`<?xml version="1.0"?><propfind xmlns="DAV:"><allprop/></propfind>`)),
		map[string]string{"Depth": "1", "Content-Type": "application/xml"},
	)

	// 5. PROPPATCH - Изменение метаданных (свойств) файла (RFC 4918)
	request("PROPPATCH", baseURL+baseFolder+"/hello.txt",
		bytes.NewBuffer([]byte(`<?xml version="1.0"?><propertyupdate xmlns="DAV:"><set><prop><description>test file</description></prop></set></propertyupdate>`)),
		map[string]string{"Content-Type": "application/xml"},
	)

	// 6. Создаем вторую папку "B"
	request("MKCOL", baseURL+workFolder+"/", nil, nil)

	// 7. COPY - Копирование файла из папки "A" в папку "B"
	request("COPY", baseURL+baseFolder+"/hello.txt",
		nil,
		map[string]string{"Destination": destURL + workFolder + "/hello_copy.txt"},
	)

	// 8. MOVE - Перемещение (переименование) файла
	request("MOVE", baseURL+workFolder+"/hello_copy.txt",
		nil,
		map[string]string{"Destination": destURL + baseFolder + "/hello_moved.txt"},
	)

	// 9. DELETE - Удаление файлов и папок
	request("DELETE", baseURL+baseFolder+"/", nil, nil)
	request("DELETE", baseURL+workFolder+"/", nil, nil)
}
