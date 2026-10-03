package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type Celebrity struct {
	ID           int    `gorm:"column:Id;primaryKey;autoIncrement:false" json:"id"`
	FullName     string `gorm:"column:FullName" json:"fullName"`
	Nationality  string `gorm:"column:Nationality" json:"nationality"`
	ReqPhotoPath string `gorm:"column:ReqPhotoPath" json:"reqPhotoPath"`
}

func (Celebrity) TableName() string {
	return "Celebrities"
}

var db *gorm.DB

const connString = "sqlserver://@LENOVO?database=CelebritiesDB&trusted_connection=true"

func main() {
	var err error
	db, err = gorm.Open(sqlserver.Open(connString), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	log.Println("Подключение к БД успешно")

	r := mux.NewRouter()
	r.HandleFunc("/Celebrities/All", getAllHandler).Methods("GET")
	r.HandleFunc("/Celebrities/{id:[0-9]+}", getByIDHandler).Methods("GET")
	r.HandleFunc("/Celebrities", createHandler).Methods("POST")
	r.HandleFunc("/Celebrities/{id:[0-9]+}", updateHandler).Methods("PUT")
	r.HandleFunc("/Celebrities/{id:[0-9]+}", deleteHandler).Methods("DELETE")
	r.Use(loggingMiddleware)

	port := ":3000"
	log.Printf("Сервер запущен на порту %s", port)
	log.Fatal(http.ListenAndServe(port, r))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func getAllHandler(w http.ResponseWriter, r *http.Request) {
	var celebrities []Celebrity
	result := db.Find(&celebrities)
	if result.Error != nil {
		log.Printf("Ошибка запроса: %v", result.Error)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(celebrities)
}

func getByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var c Celebrity
	result := db.First(&c, id)
	if result.Error == gorm.ErrRecordNotFound {
		http.Error(w, "Элемент не найден", http.StatusNotFound)
		return
	}
	if result.Error != nil {
		log.Printf("Ошибка запроса: %v", result.Error)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(c)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	var newCeleb Celebrity
	err := json.NewDecoder(r.Body).Decode(&newCeleb)
	if err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	var existing Celebrity
	result := db.First(&existing, newCeleb.ID)
	if result.Error == nil {
		http.Error(w, "Элемент с таким ID уже существует", http.StatusConflict)
		return
	}

	result = db.Create(&newCeleb)
	if result.Error != nil {
		log.Printf("Ошибка вставки: %v", result.Error)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCeleb)
	log.Printf("Добавлен элемент с ID=%d", newCeleb.ID)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var updatedCeleb Celebrity
	err := json.NewDecoder(r.Body).Decode(&updatedCeleb)
	if err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}
	updatedCeleb.ID = id

	var existing Celebrity
	result := db.First(&existing, id)
	if result.Error == gorm.ErrRecordNotFound {
		http.Error(w, "Элемент не найден", http.StatusNotFound)
		return
	}

	result = db.Model(&existing).Updates(updatedCeleb)
	if result.Error != nil {
		log.Printf("Ошибка обновления: %v", result.Error)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(updatedCeleb)
	log.Printf("Обновлен элемент с ID=%d", id)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	result := db.Delete(&Celebrity{}, id)
	if result.RowsAffected == 0 {
		http.Error(w, "Элемент не найден", http.StatusNotFound)
		return
	}
	if result.Error != nil {
		log.Printf("Ошибка удаления: %v", result.Error)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Printf("Удален элемент с ID=%d", id)
}
