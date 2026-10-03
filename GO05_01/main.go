package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Celebrity struct {
	Id           int    `json:"id"`
	FullName     string `json:"fullName"`
	Nationality  string `json:"nationality"`
	ReqPhotoPath string `json:"reqPhotoPath"`
}

const defaultDSN = "host=localhost port=5432 user=admin password=admin dbname=celebrities sslmode=disable"

var db *sql.DB

func getAllCelebrities(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /Celebrities/All")

	rows, err := db.Query("SELECT id, full_name, nationality, req_photo_path FROM celebrities ORDER BY id")
	if err != nil {
		http.Error(w, "Ошибка запроса к БД", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	celebrities := []Celebrity{}
	for rows.Next() {
		var c Celebrity
		if err := rows.Scan(&c.Id, &c.FullName, &c.Nationality, &c.ReqPhotoPath); err != nil {
			http.Error(w, "Ошибка чтения данных", http.StatusInternalServerError)
			return
		}
		celebrities = append(celebrities, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(celebrities)
}

func getCelebrityById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный Id", http.StatusBadRequest)
		return
	}
	log.Printf("GET /Celebrities/%d\n", id)

	var c Celebrity
	err = db.QueryRow(
		"SELECT id, full_name, nationality, req_photo_path FROM celebrities WHERE id = $1", id,
	).Scan(&c.Id, &c.FullName, &c.Nationality, &c.ReqPhotoPath)

	if err == sql.ErrNoRows {
		http.Error(w, "Элемент не найден", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Ошибка запроса к БД", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

func addCelebrity(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /Celebrities")

	var c Celebrity
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	var exists bool
	db.QueryRow("SELECT EXISTS(SELECT 1 FROM celebrities WHERE id = $1)", c.Id).Scan(&exists)
	if exists {
		http.Error(w, "Элемент с таким Id уже существует", http.StatusConflict) // 409
		return
	}

	if c.Id == 0 {
		err := db.QueryRow(
			"INSERT INTO celebrities (full_name, nationality, req_photo_path) VALUES ($1, $2, $3) RETURNING id",
			c.FullName, c.Nationality, c.ReqPhotoPath,
		).Scan(&c.Id)
		if err != nil {
			http.Error(w, "Ошибка вставки в БД", http.StatusInternalServerError)
			return
		}
	} else {
		_, err := db.Exec(
			"INSERT INTO celebrities (id, full_name, nationality, req_photo_path) VALUES ($1, $2, $3, $4)",
			c.Id, c.FullName, c.Nationality, c.ReqPhotoPath,
		)
		if err != nil {
			http.Error(w, "Ошибка вставки в БД", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

func updateCelebrity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный Id", http.StatusBadRequest)
		return
	}
	log.Printf("PUT /Celebrities/%d\n", id)

	var c Celebrity
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	result, err := db.Exec(
		"UPDATE celebrities SET full_name=$1, nationality=$2, req_photo_path=$3 WHERE id=$4",
		c.FullName, c.Nationality, c.ReqPhotoPath, id,
	)
	if err != nil {
		http.Error(w, "Ошибка обновления в БД", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Элемент не найден", http.StatusNotFound) // 404
		return
	}

	c.Id = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(c)
}

func deleteCelebrity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный Id", http.StatusBadRequest)
		return
	}
	log.Printf("DELETE /Celebrities/%d\n", id)

	result, err := db.Exec("DELETE FROM celebrities WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Ошибка удаления из БД", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Элемент не найден", http.StatusNotFound) // 404
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Элемент удалён"))
}

func main() {
	dsn := defaultDSN
	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("БД недоступна:", err)
	}
	log.Println("Подключение к PostgreSQL успешно")

	router := mux.NewRouter()
	router.HandleFunc("/Celebrities/All", getAllCelebrities).Methods("GET")
	router.HandleFunc("/Celebrities/{id}", getCelebrityById).Methods("GET")
	router.HandleFunc("/Celebrities", addCelebrity).Methods("POST")
	router.HandleFunc("/Celebrities/{id}", updateCelebrity).Methods("PUT")
	router.HandleFunc("/Celebrities/{id}", deleteCelebrity).Methods("DELETE")

	log.Println("Сервер запущен на порту 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
