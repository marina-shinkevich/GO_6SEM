package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)


type Celebrity struct {
	ID           int    `json:"id" gorm:"primaryKey;column:Id"`
	FullName     string `json:"fullName" gorm:"column:FullName"`
	Nationality  string `json:"nationality" gorm:"column:Nationality"`
	ReqPhotoPath string `json:"reqPhotoPath" gorm:"column:ReqPhotoPath"`
}


func (Celebrity) TableName() string {
	return "Celebrities"
}

var db *gorm.DB

func main() {
	
	dsn := "sqlserver://celebrities_app:AppPassword123!@localhost:1433?database=CelebritiesDB"
	var err error
	
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
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}
	log.Println("Подключение к базе данных установлено")


	
	router := mux.NewRouter()

	router.HandleFunc("/Celebrities/All", getAllCelebrities).Methods("GET")
	router.HandleFunc("/Celebrities/{id}", getCelebrityByID).Methods("GET")
	router.HandleFunc("/Celebrities", addCelebrity).Methods("POST")
	router.HandleFunc("/Celebrities/{id}", updateCelebrity).Methods("PUT")
	router.HandleFunc("/Celebrities/{id}", deleteCelebrity).Methods("DELETE")


	log.Println("Сервер запущен на порту 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}

func getAllCelebrities(w http.ResponseWriter, r *http.Request) {
	var celebrities []Celebrity
	result := db.Find(&celebrities)
	if result.Error != nil {
		http.Error(w, "Ошибка при получении данных: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(celebrities)
	log.Println("GET /Celebrities/All - возвращено записей:", len(celebrities))
}

func getCelebrityByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var celebrity Celebrity
	result := db.First(&celebrity, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "Запись не найдена", http.StatusNotFound)
		} else {
			http.Error(w, "Ошибка при получении данных: "+result.Error.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(celebrity)
	log.Printf("GET /Celebrities/%d - запись найдена\n", id)
}


func addCelebrity(w http.ResponseWriter, r *http.Request) {
	var celebrity Celebrity
	err := json.NewDecoder(r.Body).Decode(&celebrity)
	if err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	var existing Celebrity
	result := db.First(&existing, celebrity.ID)
	if result.Error == nil {
		http.Error(w, "Запись с таким ID уже существует", http.StatusConflict)
		log.Printf("POST /Celebrities - конфликт: ID %d уже существует\n", celebrity.ID)
		return
	}


	
	result = db.Exec("INSERT INTO Celebrities (Id, FullName, Nationality, ReqPhotoPath) VALUES (?, ?, ?, ?)",
		celebrity.ID, celebrity.FullName, celebrity.Nationality, celebrity.ReqPhotoPath)
	if result.Error != nil {
		http.Error(w, "Ошибка при добавлении записи: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(celebrity)
	log.Printf("POST /Celebrities - добавлена запись с ID %d\n", celebrity.ID)
}


func updateCelebrity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var celebrity Celebrity
	err = json.NewDecoder(r.Body).Decode(&celebrity)
	if err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}


	var existing Celebrity
	result := db.First(&existing, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "Запись не найдена", http.StatusNotFound)
			log.Printf("PUT /Celebrities/%d - запись не найдена\n", id)
		} else {
			http.Error(w, "Ошибка при проверке записи: "+result.Error.Error(), http.StatusInternalServerError)
		}
		return
	}

	celebrity.ID = id
	result = db.Save(&celebrity)
	if result.Error != nil {
		http.Error(w, "Ошибка при обновлении записи: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(celebrity)
	log.Printf("PUT /Celebrities/%d - запись обновлена\n", id)
}


func deleteCelebrity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}


	var existing Celebrity
	result := db.First(&existing, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "Запись не найдена", http.StatusNotFound)
			log.Printf("DELETE /Celebrities/%d - запись не найдена\n", id)
		} else {
			http.Error(w, "Ошибка при проверке записи: "+result.Error.Error(), http.StatusInternalServerError)
		}
		return
	}


	result = db.Delete(&Celebrity{}, id)
	if result.Error != nil {
		http.Error(w, "Ошибка при удалении записи: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Printf("DELETE /Celebrities/%d - запись удалена\n", id)
}