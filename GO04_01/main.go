package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
)

type Celebrity struct {
	Id           int    `json:"id"`
	FullName     string `json:"fullName"`
	Nationality  string `json:"nationality"`
	ReqPhotoPath string `json:"reqPhotoPath"`
}

func dataFile() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal("Не удалось определить путь к исполняемому файлу:", err)
	}

	p := filepath.Join(filepath.Dir(exePath), "data", "Celebrities.json")
	if _, err := os.Stat(p); err == nil {
		return p
	}

	return filepath.Join("data", "Celebrities.json")
}

func loadCelebrities() ([]Celebrity, error) {
	data, err := os.ReadFile(dataFile())
	if err != nil {
		return nil, err
	}
	var celebrities []Celebrity
	err = json.Unmarshal(data, &celebrities)
	return celebrities, err
}

func saveCelebrities(celebrities []Celebrity) error {
	data, err := json.MarshalIndent(celebrities, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile(), data, 0644)
}

func getAllCelebrities(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /Celebrities/All")

	celebrities, err := loadCelebrities()
	if err != nil {
		http.Error(w, "Ошибка чтения файла", http.StatusInternalServerError)
		return
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

	celebrities, err := loadCelebrities()
	if err != nil {
		http.Error(w, "Ошибка чтения файла", http.StatusInternalServerError)
		return
	}

	for _, c := range celebrities {
		if c.Id == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(c)
			return
		}
	}

	http.Error(w, "Элемент не найден", http.StatusNotFound)
}

func addCelebrity(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /Celebrities")

	var newCelebrity Celebrity
	err := json.NewDecoder(r.Body).Decode(&newCelebrity)
	if err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	celebrities, err := loadCelebrities()
	if err != nil {
		http.Error(w, "Ошибка чтения файла", http.StatusInternalServerError)
		return
	}

	for _, c := range celebrities {
		if c.Id == newCelebrity.Id {
			http.Error(w, "Элемент с таким Id уже существует", http.StatusConflict) // 409
			return
		}
	}

	celebrities = append(celebrities, newCelebrity)
	if err := saveCelebrities(celebrities); err != nil {
		http.Error(w, "Ошибка сохранения файла", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCelebrity)
}

func updateCelebrity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный Id", http.StatusBadRequest)
		return
	}
	log.Printf("PUT /Celebrities/%d\n", id)

	var updated Celebrity
	err = json.NewDecoder(r.Body).Decode(&updated)
	if err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	celebrities, err := loadCelebrities()
	if err != nil {
		http.Error(w, "Ошибка чтения файла", http.StatusInternalServerError)
		return
	}

	for i, c := range celebrities {
		if c.Id == id {
			updated.Id = id
			celebrities[i] = updated
			if err := saveCelebrities(celebrities); err != nil {
				http.Error(w, "Ошибка сохранения файла", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updated)
			return
		}
	}

	http.Error(w, "Элемент не найден", http.StatusNotFound) // 404
}

func deleteCelebrity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный Id", http.StatusBadRequest)
		return
	}
	log.Printf("DELETE /Celebrities/%d\n", id)

	celebrities, err := loadCelebrities()
	if err != nil {
		http.Error(w, "Ошибка чтения файла", http.StatusInternalServerError)
		return
	}

	for i, c := range celebrities {
		if c.Id == id {
			celebrities = append(celebrities[:i], celebrities[i+1:]...)
			if err := saveCelebrities(celebrities); err != nil {
				http.Error(w, "Ошибка сохранения файла", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Элемент удалён"))
			return
		}
	}

	http.Error(w, "Элемент не найден", http.StatusNotFound) // 404
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/Celebrities/All", getAllCelebrities).Methods("GET")
	router.HandleFunc("/Celebrities/{id}", getCelebrityById).Methods("GET")
	router.HandleFunc("/Celebrities", addCelebrity).Methods("POST")
	router.HandleFunc("/Celebrities/{id}", updateCelebrity).Methods("PUT")
	router.HandleFunc("/Celebrities/{id}", deleteCelebrity).Methods("DELETE")

	log.Println("Сервер запущен на порту 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
