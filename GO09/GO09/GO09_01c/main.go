package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	baseURL  = "http://localhost:8082"
	username = "webdavuser"
	password = "mypassword123"
)

func doRequest(method, path, contentType string, body io.Reader) (*http.Response, error) {
	url := baseURL + "/" + strings.TrimPrefix(path, "/")
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(username, password)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	return http.DefaultClient.Do(req)
}

// MKCOL
func mkcol(path string) {
	urlPath := strings.TrimRight(path, "/") + "/"
	resp, err := doRequest("MKCOL", urlPath, "", http.NoBody)
	if err != nil {
		fmt.Printf("[MKCOL] Ошибка: %v\n", err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("[MKCOL] %s -> %s\n", path, resp.Status)
}

// PUT
func put(remotePath, localFile string) {
	f, err := os.Open(localFile)
	if err != nil {
		fmt.Printf("[PUT] Не удалось открыть файл %s: %v\n", localFile, err)
		return
	}
	defer f.Close()

	resp, err := doRequest("PUT", remotePath, "application/octet-stream", f)
	if err != nil {
		fmt.Printf("[PUT] Ошибка: %v\n", err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("[PUT] %s -> %s\n", remotePath, resp.Status)
}

// GET
func get(remotePath, localFile string) {
	resp, err := doRequest("GET", remotePath, "", nil)
	if err != nil {
		fmt.Printf("[GET] Ошибка: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("[GET] %s -> %s\n", remotePath, resp.Status)
		return
	}

	out, err := os.Create(localFile)
	if err != nil {
		fmt.Printf("[GET] Не удалось создать файл %s: %v\n", localFile, err)
		return
	}
	defer out.Close()

	io.Copy(out, resp.Body)
	fmt.Printf("[GET] %s сохранён в %s\n", remotePath, localFile)
}

// COPY
func copyRes(src, dst string) {
	url := baseURL + "/" + strings.TrimPrefix(src, "/")
	dstURL := baseURL + "/" + strings.TrimPrefix(dst, "/")

	req, err := http.NewRequest("COPY", url, nil)
	if err != nil {
		fmt.Printf("[COPY] Ошибка: %v\n", err)
		return
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Destination", dstURL)
	req.Header.Set("Overwrite", "T")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("[COPY] Ошибка: %v\n", err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("[COPY] %s -> %s : %s\n", src, dst, resp.Status)
}

// MOVE
func moveRes(src, dst string) {
	url := baseURL + "/" + strings.TrimPrefix(src, "/")
	dstURL := baseURL + "/" + strings.TrimPrefix(dst, "/")

	req, err := http.NewRequest("MOVE", url, nil)
	if err != nil {
		fmt.Printf("[MOVE] Ошибка: %v\n", err)
		return
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Destination", dstURL)
	req.Header.Set("Overwrite", "T")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("[MOVE] Ошибка: %v\n", err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("[MOVE] %s -> %s : %s\n", src, dst, resp.Status)
}

// DELETE
func deleteRes(path string) {
	resp, err := doRequest("DELETE", path, "", nil)
	if err != nil {
		fmt.Printf("[DELETE] Ошибка: %v\n", err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("[DELETE] %s -> %s\n", path, resp.Status)
}

// RMCOL
func rmcol(path string) {
	urlPath := strings.TrimRight(path, "/") + "/"
	url := baseURL + "/" + strings.TrimPrefix(urlPath, "/")
	req, err := http.NewRequest("DELETE", url, http.NoBody)
	if err != nil {
		fmt.Printf("[RMCOL] Ошибка: %v\n", err)
		return
	}
	req.SetBasicAuth(username, password)
	req.Header.Set("Depth", "infinity")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("[RMCOL] Ошибка: %v\n", err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf("[RMCOL] %s -> %s\n", path, resp.Status)
}

func main() {
	fmt.Println("=== WebDAV Client GO09_01c ===")

	fmt.Println("\n--- MKCOL ---")
	mkcol("/testdir")

	fmt.Println("\n--- PUT ---")
	os.WriteFile("test.txt", []byte("Привет, WebDAV!\n"), 0644)
	put("/testdir/test.txt", "test.txt")

	fmt.Println("\n--- GET ---")
	get("/testdir/test.txt", "downloaded.txt")
	data, _ := os.ReadFile("downloaded.txt")
	fmt.Printf("Содержимое: %s", string(data))

	fmt.Println("\n--- COPY ---")
	copyRes("/testdir/test.txt", "/testdir/test_copy.txt")

	fmt.Println("\n--- MOVE ---")
	moveRes("/testdir/test_copy.txt", "/testdir/test_moved.txt")

	fmt.Println("\n--- DELETE (файлы) ---")
	deleteRes("/testdir/test.txt")
	deleteRes("/testdir/test_moved.txt")

	fmt.Println("\n--- DELETE (директория) ---")
	rmcol("/testdir")

	fmt.Println("\n=== Готово ===")
}
