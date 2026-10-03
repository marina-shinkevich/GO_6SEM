package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "GO11_01/docs"

	"github.com/gorilla/mux"
	_ "github.com/microsoft/go-mssqldb"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Celebrity struct {
	ID           int    `json:"id" example:"1"`
	FullName     string `json:"fullName" example:"Charlie Chaplin"`
	Nationality  string `json:"nationality" example:"United Kingdom"`
	ReqPhotoPath string `json:"reqPhotoPath" example:"photos/charlie_chaplin.jpg"`
}

const (
	defaultSQLHost  = `localhost`
	defaultSQLPort  = ``
	defaultDatabase = `PIS_Lab5`
)

var db *sql.DB

// @title GO11_01 API
// @version 1.0
// @description REST-интерфейс для работы с коллекцией Celebrities.
// @host localhost:3000
// @BasePath /
func main() {
	connString, err := openAppDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	initDB()

	r := mux.NewRouter()
	r.Use(logMiddleware)
	r.HandleFunc("/", index).Methods("GET")
	r.HandleFunc("/Celebrities/All", getAll).Methods("GET")
	r.HandleFunc("/Celebrities/{id:[0-9]+}", getOne).Methods("GET")
	r.HandleFunc("/Celebrities", addOne).Methods("POST")
	r.HandleFunc("/Celebrities/{id:[0-9]+}", updateOne).Methods("PUT")
	r.HandleFunc("/Celebrities/{id:[0-9]+}", deleteOne).Methods("DELETE")
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("cannot start GO11_01: port 3000 is already busy (%v). Close the running server or run: netstat -ano | findstr :3000", err)
	}

	log.Printf("GO11_01 started on http://localhost:3000 using %s", connString)
	log.Printf("Swagger UI: http://localhost:3000/swagger/index.html")
	log.Fatal(http.Serve(listener, r))
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("GO11_01 REST API. Swagger UI: /swagger/index.html\n"))
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func initDB() {
	_, err := db.Exec(`
		IF OBJECT_ID(N'dbo.celebrities', N'U') IS NULL
		BEGIN
			CREATE TABLE dbo.celebrities (
				id INT NOT NULL PRIMARY KEY,
				full_name NVARCHAR(200) NOT NULL,
				nationality NVARCHAR(120) NOT NULL,
				req_photo_path NVARCHAR(260) NOT NULL
			)
		END
	`)
	if err != nil {
		log.Fatal(err)
	}

	var count int
	if err := db.QueryRow(`SELECT COUNT(*) FROM dbo.celebrities`).Scan(&count); err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		_, err = db.Exec(`
			INSERT INTO dbo.celebrities (id, full_name, nationality, req_photo_path) VALUES
			(1, N'Charlie Chaplin', N'United Kingdom', N'photos/charlie_chaplin.jpg'),
			(2, N'Audrey Hepburn', N'Belgium', N'photos/audrey_hepburn.jpg'),
			(3, N'Bruce Lee', N'China', N'photos/bruce_lee.jpg')
		`)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// getAll godoc
// @Summary Вся коллекция
// @Description Получить все элементы коллекции Celebrities.
// @Tags Celebrities
// @Produce json
// @Success 200 {array} Celebrity "Успешный ответ"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /Celebrities/All [get]
func getAll(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`SELECT id, full_name, nationality, req_photo_path FROM dbo.celebrities ORDER BY id`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []Celebrity
	for rows.Next() {
		var c Celebrity
		if err := rows.Scan(&c.ID, &c.FullName, &c.Nationality, &c.ReqPhotoPath); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		items = append(items, c)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, items)
}

// getOne godoc
// @Summary Элемент коллекции
// @Description Получить элемент коллекции Celebrities по id.
// @Tags Celebrities
// @Produce json
// @Param id path int true "Id элемента"
// @Success 200 {object} Celebrity "Успешный ответ"
// @Failure 404 {string} string "Элемент не найден"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /Celebrities/{id} [get]
func getOne(w http.ResponseWriter, r *http.Request) {
	id := readID(r)
	var c Celebrity

	err := db.QueryRow(
		`SELECT id, full_name, nationality, req_photo_path FROM dbo.celebrities WHERE id = @p1`,
		id,
	).Scan(&c.ID, &c.FullName, &c.Nationality, &c.ReqPhotoPath)

	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, c)
}

// addOne godoc
// @Summary Добавить элемент
// @Description Добавить элемент в коллекцию Celebrities. При дублировании id возвращается 409.
// @Tags Celebrities
// @Accept json
// @Produce json
// @Param celebrity body Celebrity true "Данные элемента"
// @Success 201 {object} Celebrity "Элемент создан"
// @Failure 400 {string} string "Некорректный JSON"
// @Failure 409 {string} string "Элемент с таким id уже существует"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /Celebrities [post]
func addOne(w http.ResponseWriter, r *http.Request) {
	var c Celebrity
	if json.NewDecoder(r.Body).Decode(&c) != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	var id int
	err := db.QueryRow(`SELECT id FROM dbo.celebrities WHERE id = @p1`, c.ID).Scan(&id)
	if err == nil {
		http.Error(w, "id already exists", http.StatusConflict)
		return
	}
	if !errors.Is(err, sql.ErrNoRows) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.Exec(
		`INSERT INTO dbo.celebrities (id, full_name, nationality, req_photo_path) VALUES (@p1, @p2, @p3, @p4)`,
		c.ID, c.FullName, c.Nationality, c.ReqPhotoPath,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	writeJSON(w, c)
}

// updateOne godoc
// @Summary Изменить элемент
// @Description Изменить элемент коллекции Celebrities по id. При отсутствии элемента возвращается 404.
// @Tags Celebrities
// @Accept json
// @Produce json
// @Param id path int true "Id элемента"
// @Param celebrity body Celebrity true "Данные элемента"
// @Success 200 {object} Celebrity "Успешный ответ"
// @Failure 400 {string} string "Некорректный JSON"
// @Failure 404 {string} string "Элемент не найден"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /Celebrities/{id} [put]
func updateOne(w http.ResponseWriter, r *http.Request) {
	id := readID(r)
	var c Celebrity
	if json.NewDecoder(r.Body).Decode(&c) != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	c.ID = id

	res, err := db.Exec(
		`UPDATE dbo.celebrities SET full_name = @p1, nationality = @p2, req_photo_path = @p3 WHERE id = @p4`,
		c.FullName, c.Nationality, c.ReqPhotoPath, id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	n, _ := res.RowsAffected()
	if n == 0 {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	writeJSON(w, c)
}

// deleteOne godoc
// @Summary Удалить элемент
// @Description Удалить элемент коллекции Celebrities по id. При отсутствии элемента возвращается 404.
// @Tags Celebrities
// @Produce json
// @Param id path int true "Id элемента"
// @Success 200 {object} Celebrity "Успешный ответ"
// @Failure 404 {string} string "Элемент не найден"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /Celebrities/{id} [delete]
func deleteOne(w http.ResponseWriter, r *http.Request) {
	id := readID(r)
	var c Celebrity

	err := db.QueryRow(
		`SELECT id, full_name, nationality, req_photo_path FROM dbo.celebrities WHERE id = @p1`,
		id,
	).Scan(&c.ID, &c.FullName, &c.Nationality, &c.ReqPhotoPath)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := db.Exec(`DELETE FROM dbo.celebrities WHERE id = @p1`, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, c)
}

func readID(r *http.Request) int {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	return id
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(v)
}

func openAppDB() (string, error) {
	candidates := connectionCandidates(databaseName())
	var errs []string

	for _, connString := range candidates {
		currentDB, err := openAndInitDatabase(connString)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%s -> %v", connString, err))
			continue
		}

		db = currentDB
		return connString, nil
	}

	return "", fmt.Errorf("sql server connection failed:\n%s", strings.Join(errs, "\n"))
}

func openAndInitDatabase(connString string) (*sql.DB, error) {
	masterDB, err := sql.Open("sqlserver", strings.Replace(connString, "database="+databaseName(), "database=master", 1))
	if err != nil {
		return nil, err
	}
	defer masterDB.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := masterDB.PingContext(ctx); err != nil {
		return nil, err
	}

	dbName := databaseName()
	if _, err := masterDB.ExecContext(ctx, fmt.Sprintf(`
		IF DB_ID(N'%s') IS NULL
		BEGIN
			CREATE DATABASE [%s]
		END
	`, escapeSQLString(dbName), escapeSQLIdentifier(dbName))); err != nil {
		return nil, err
	}

	appDB, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, err
	}

	if err := appDB.PingContext(ctx); err != nil {
		appDB.Close()
		return nil, err
	}

	return appDB, nil
}

func connectionCandidates(dbName string) []string {
	if raw := strings.TrimSpace(os.Getenv("GO11_CONNECTION_STRING")); raw != "" {
		return []string{raw}
	}

	server := envOrDefault("GO11_SQLSERVER", defaultSQLHost)
	port := envOrDefault("GO11_SQLPORT", defaultSQLPort)
	base := []string{
		buildConnectionString(server, port, dbName),
		buildConnectionString(`.\`, ``, dbName),           // named pipes для default instance
		buildConnectionString(`localhost\SQLEXPRESS`, ``, dbName), // SQL Express
		buildConnectionString(`.\SQLEXPRESS`, ``, dbName), // SQL Express через named pipes
	}

	return uniqueStrings(base)
}

func buildConnectionString(server, port, dbName string) string {
	base := []string{
		fmt.Sprintf("server=%s", server),
		fmt.Sprintf("database=%s", dbName),
		"encrypt=true",
		"TrustServerCertificate=true",
		"connection timeout=5",
		"dial timeout=3",
	}

	if strings.TrimSpace(port) != "" {
		base = append(base[:1], append([]string{fmt.Sprintf("port=%s", port)}, base[1:]...)...)
	}

	return strings.Join(base, ";")
}

func databaseName() string {
	return envOrDefault("GO11_DATABASE", defaultDatabase)
}

func envOrDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func escapeSQLString(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}

func escapeSQLIdentifier(value string) string {
	return strings.ReplaceAll(value, "]", "]]")
}

func uniqueStrings(items []string) []string {
	seen := make(map[string]struct{}, len(items))
	result := make([]string, 0, len(items))
	for _, item := range items {
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		result = append(result, item)
	}
	return result
}
