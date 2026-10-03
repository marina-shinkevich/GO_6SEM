package main

import (
	"GO03_02/P03_02"
	"fmt"
	"log"
	"net/http"
)

var stats = &P03_02.RequestStats{}

func sHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		stats.PlusGet()
		fmt.Fprint(w, "GET /S принят")
	case http.MethodPost:
		stats.PlusPost()
		fmt.Fprint(w, "POST /S принят")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func gHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, stats.GenStr())
	}
}

func main() {
	http.HandleFunc("/S", sHandler)
	http.HandleFunc("/G", gHandler)

	log.Println("Сервер GO03_02 запущен на порту :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
