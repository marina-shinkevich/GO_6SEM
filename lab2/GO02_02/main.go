package main

import (
	"fmt"
	"go02_02/go02_02lib"
	"net/http"
)

var A01 = 3

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "A01 = %d\n", A01)
	fmt.Fprintf(w, "A02 = %t\n", A02)
	fmt.Fprintf(w, "A03 = %s\n", go02_02lib.A03)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server started on port 4000")
	http.ListenAndServe(":4000", nil)
}
