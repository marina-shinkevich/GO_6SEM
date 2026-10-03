package main

import (
	"fmt"
	"net/http"

	"GO02_01/go02_01lib"
)

const C01 = 3.14

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintf(w, "C01 = %g\n", C01)
	fmt.Fprintf(w, "C02 = %e\n", C02)
	fmt.Fprintf(w, "C03 = %g\n", go02_01lib.C03)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("Server started on port 3000")
	http.ListenAndServe(":3000", nil)
}
