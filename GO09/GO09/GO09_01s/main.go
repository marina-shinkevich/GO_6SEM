package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	storageDir = "./webdav"
	serverPort = ":8082"
	username   = "webdavuser"
	password   = "mypassword123"
)

func realPath(r *http.Request) string {
	clean := filepath.Clean(r.URL.Path)
	clean = strings.TrimPrefix(clean, "/")
	return filepath.Join(storageDir, clean)
}

func checkAuth(r *http.Request) bool {
	u, p, ok := r.BasicAuth()
	return ok && u == username && p == password
}

func handler(w http.ResponseWriter, r *http.Request) {
	if !checkAuth(r) {
		w.Header().Set("WWW-Authenticate", `Basic realm="WebDAV Server"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		fmt.Printf("[%s] %s -> 401\n", r.Method, r.URL.Path)
		return
	}

	switch r.Method {
	case "MKCOL":
		handleMkcol(w, r)
	case "PUT":
		handlePut(w, r)
	case "GET":
		handleGet(w, r)
	case "COPY":
		handleCopy(w, r)
	case "MOVE":
		handleMove(w, r)
	case "DELETE":
		handleDelete(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		fmt.Printf("[%s] %s -> 405\n", r.Method, r.URL.Path)
	}
}

func handleMkcol(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimRight(realPath(r), "/")
	if _, err := os.Stat(path); err == nil {
		http.Error(w, "Collection already exists", http.StatusMethodNotAllowed)
		fmt.Printf("[MKCOL] %s -> 405\n", r.URL.Path)
		return
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		fmt.Printf("[MKCOL] %s -> 409\n", r.URL.Path)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Printf("[MKCOL] %s -> 201\n", r.URL.Path)
}

func handlePut(w http.ResponseWriter, r *http.Request) {
	path := realPath(r)
	os.MkdirAll(filepath.Dir(path), 0755)
	f, err := os.Create(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Printf("[PUT] %s -> 500\n", r.URL.Path)
		return
	}
	defer f.Close()
	io.Copy(f, r.Body)
	w.WriteHeader(http.StatusCreated)
	fmt.Printf("[PUT] %s -> 201\n", r.URL.Path)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	path := realPath(r)
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		http.Error(w, "Not found", http.StatusNotFound)
		fmt.Printf("[GET] %s -> 404\n", r.URL.Path)
		return
	}
	if info.IsDir() {
		http.Error(w, "Is a directory", http.StatusBadRequest)
		fmt.Printf("[GET] %s -> 400\n", r.URL.Path)
		return
	}
	http.ServeFile(w, r, path)
	fmt.Printf("[GET] %s -> 200\n", r.URL.Path)
}

func handleCopy(w http.ResponseWriter, r *http.Request) {
	srcPath := realPath(r)
	dstURL := r.Header.Get("Destination")
	if dstURL == "" {
		http.Error(w, "Destination header required", http.StatusBadRequest)
		fmt.Printf("[COPY] %s -> 400\n", r.URL.Path)
		return
	}
	dstParsed, _ := url.Parse(dstURL)
	dstPath := realPath(&http.Request{URL: dstParsed})

	src, err := os.Open(srcPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		fmt.Printf("[COPY] %s -> 404\n", r.URL.Path)
		return
	}
	defer src.Close()

	os.MkdirAll(filepath.Dir(dstPath), 0755)
	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Printf("[COPY] %s -> 500\n", r.URL.Path)
		return
	}
	defer dst.Close()

	io.Copy(dst, src)
	w.WriteHeader(http.StatusNoContent)
	fmt.Printf("[COPY] %s -> 204\n", r.URL.Path)
}

func handleMove(w http.ResponseWriter, r *http.Request) {
	srcPath := realPath(r)
	dstURL := r.Header.Get("Destination")
	if dstURL == "" {
		http.Error(w, "Destination header required", http.StatusBadRequest)
		fmt.Printf("[MOVE] %s -> 400\n", r.URL.Path)
		return
	}
	dstParsed, _ := url.Parse(dstURL)
	dstPath := realPath(&http.Request{URL: dstParsed})

	os.MkdirAll(filepath.Dir(dstPath), 0755)
	if err := os.Rename(srcPath, dstPath); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Printf("[MOVE] %s -> 500\n", r.URL.Path)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	fmt.Printf("[MOVE] %s -> 204\n", r.URL.Path)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	path := realPath(r)
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		http.Error(w, "Not found", http.StatusNotFound)
		fmt.Printf("[DELETE] %s -> 404\n", r.URL.Path)
		return
	}
	if info.IsDir() && r.Header.Get("Depth") != "infinity" {
		http.Error(w, "Is a directory, use Depth: infinity", http.StatusConflict)
		fmt.Printf("[DELETE] %s -> 409\n", r.URL.Path)
		return
	}
	if info.IsDir() {
		err = os.RemoveAll(path)
	} else {
		err = os.Remove(path)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Printf("[DELETE] %s -> 500\n", r.URL.Path)
		return
	}
	w.WriteHeader(http.StatusNoContent)
	fmt.Printf("[DELETE] %s -> 204\n", r.URL.Path)
}

func main() {
	os.MkdirAll(storageDir, 0755)
	fmt.Printf("WebDAV server starting on %s\n", serverPort)
	if err := http.ListenAndServe(serverPort, http.HandlerFunc(handler)); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}
