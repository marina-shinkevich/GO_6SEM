package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Input struct {
	Op string `json:"op"`
	X  int    `json:"x"`
	Y  int    `json:"y"`
}

type Output struct {
	Op     string `json:"op"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Result int    `json:"result"`
}

func handlerJSON(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var in Input

	
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
	
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	var result int
	var isResult bool = true

	switch in.Op {
	case "+":
		result = in.X + in.Y
	case "-":
		result = in.X - in.Y
	case "*":
		result = in.X * in.Y
	default:
		isResult = false
	}

	if isResult {

		out := Output{Op: in.Op, X: in.X, Y: in.Y, Result: result}

	
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		if err := json.NewEncoder(w).Encode(out); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "unsupported operation", http.StatusBadRequest)
	}
}

func main() {
	http.HandleFunc("/JSON", handlerJSON)
	log.Println("Сервер запущен: http://localhost:3000/JSON")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
