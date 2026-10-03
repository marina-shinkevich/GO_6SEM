package main

import (
	"log"
	"net/http"

	"golang.org/x/net/webdav" // Подключаем пакет WebDAV
)

func main() {
	// 1. Инициализируем WebDAV-обработчик
	handler := &webdav.Handler{
		// Указываем корневой префикс URL
		Prefix: "/",

		// Задаем физическую папку на диске сервера, куда будут сохраняться файлы.
		// Папка "./data" должна быть создана заранее!
		FileSystem: webdav.Dir("./data"),
	}

	log.Println("WebDAV server on :3000")

	// 2. Запускаем сервер, передавая наш WebDAV-обработчик
	err := http.ListenAndServe(":3000", handler)
	if err != nil {
		log.Fatal(err)
	}
}
