package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	_ "github.com/microsoft/go-mssqldb"
	_ "github.com/microsoft/go-mssqldb/sharedmemory"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Celebrity struct {
	ID           int    `json:"id" gorm:"column:id;primaryKey;autoIncrement:false"`
	FullName     string `json:"fullName" gorm:"column:full_name"`
	Nationality  string `json:"nationality" gorm:"column:nationality"`
	ReqPhotoPath string `json:"reqPhotoPath" gorm:"column:req_photo_path"`
}

func (Celebrity) TableName() string {
	return "celebrities"
}

const (
	defaultSQLHost  = `localhost`
	defaultSQLPort  = ``
	defaultDatabase = `PIS_Lab5`
)

const gqlSchema = `
	schema {
		query: Query
		mutation: Mutation
	}

	type Query {
		celebrities: [Celebrity!]!
		celebrity(id: Int!): Celebrity
	}

	type Mutation {
		addCelebrity(input: AddCelebrityInput!): Celebrity!
		updateCelebrity(id: Int!, input: UpdateCelebrityInput!): Celebrity!
		deleteCelebrity(id: Int!): Celebrity!
	}

	type Celebrity {
		id: Int!
		fullName: String!
		nationality: String!
		reqPhotoPath: String!
	}

	input AddCelebrityInput {
		id: Int!
		fullName: String!
		nationality: String!
		reqPhotoPath: String!
	}

	input UpdateCelebrityInput {
		fullName: String!
		nationality: String!
		reqPhotoPath: String!
	}
`

var db *gorm.DB

func main() {
	connString, err := openAppDB()
	if err != nil {
		log.Fatal(err)
	}

	initDB()

	schema := graphql.MustParseSchema(gqlSchema, &resolver{})
	http.Handle("/", http.HandlerFunc(index))
	http.Handle("/graphql", logMiddleware(&relay.Handler{Schema: schema}))

	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("cannot start GO10_01: port 3000 is already busy (%v). Close the running server or run: netstat -ano | findstr :3000", err)
	}

	log.Printf("GO10_01 started on http://localhost:3000/graphql using %s", connString)
	log.Fatal(http.Serve(listener, nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("GO10_01 GraphQL endpoint: POST /graphql\n"))
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func initDB() {
	if err := db.AutoMigrate(&Celebrity{}); err != nil {
		log.Fatal(err)
	}

	var count int64
	if err := db.Model(&Celebrity{}).Count(&count).Error; err != nil {
		log.Fatal(err)
	}
	if count == 0 {
		if err := db.Create([]Celebrity{
			{ID: 1, FullName: "Charlie Chaplin", Nationality: "United Kingdom", ReqPhotoPath: "photos/charlie_chaplin.jpg"},
			{ID: 2, FullName: "Audrey Hepburn", Nationality: "Belgium", ReqPhotoPath: "photos/audrey_hepburn.jpg"},
			{ID: 3, FullName: "Bruce Lee", Nationality: "China", ReqPhotoPath: "photos/bruce_lee.jpg"},
		}).Error; err != nil {
			log.Fatal(err)
		}
	}
}

type resolver struct{}

type idArgs struct {
	ID int32
}

type addCelebrityArgs struct {
	Input addCelebrityInput
}

type addCelebrityInput struct {
	ID           int32
	FullName     string
	Nationality  string
	ReqPhotoPath string
}

type updateCelebrityArgs struct {
	ID    int32
	Input updateCelebrityInput
}

type updateCelebrityInput struct {
	FullName     string
	Nationality  string
	ReqPhotoPath string
}

func (r *resolver) Celebrities() ([]*celebrityResolver, error) {
	var items []Celebrity
	if err := db.Order("id").Find(&items).Error; err != nil {
		return nil, err
	}

	result := make([]*celebrityResolver, 0, len(items))
	for i := range items {
		result = append(result, &celebrityResolver{c: items[i]})
	}
	return result, nil
}

func (r *resolver) Celebrity(args idArgs) (*celebrityResolver, error) {
	var c Celebrity
	err := db.First(&c, int(args.ID)).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &celebrityResolver{c: c}, nil
}

func (r *resolver) AddCelebrity(args addCelebrityArgs) (*celebrityResolver, error) {
	c := Celebrity{
		ID:           int(args.Input.ID),
		FullName:     args.Input.FullName,
		Nationality:  args.Input.Nationality,
		ReqPhotoPath: args.Input.ReqPhotoPath,
	}

	var old Celebrity
	err := db.First(&old, c.ID).Error
	if err == nil {
		return nil, fmt.Errorf("id already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if err := db.Create(&c).Error; err != nil {
		return nil, err
	}
	return &celebrityResolver{c: c}, nil
}

func (r *resolver) UpdateCelebrity(args updateCelebrityArgs) (*celebrityResolver, error) {
	id := int(args.ID)
	updates := map[string]any{
		"full_name":      args.Input.FullName,
		"nationality":    args.Input.Nationality,
		"req_photo_path": args.Input.ReqPhotoPath,
	}

	res := db.Model(&Celebrity{}).Where("id = ?", id).Updates(updates)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, fmt.Errorf("not found")
	}

	c := Celebrity{
		ID:           id,
		FullName:     args.Input.FullName,
		Nationality:  args.Input.Nationality,
		ReqPhotoPath: args.Input.ReqPhotoPath,
	}
	return &celebrityResolver{c: c}, nil
}

func (r *resolver) DeleteCelebrity(args idArgs) (*celebrityResolver, error) {
	id := int(args.ID)
	var c Celebrity

	err := db.First(&c, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("not found")
	}
	if err != nil {
		return nil, err
	}

	if err := db.Delete(&Celebrity{}, id).Error; err != nil {
		return nil, err
	}
	return &celebrityResolver{c: c}, nil
}

type celebrityResolver struct {
	c Celebrity
}

func (r *celebrityResolver) ID() int32 {
	return int32(r.c.ID)
}

func (r *celebrityResolver) FullName() string {
	return r.c.FullName
}

func (r *celebrityResolver) Nationality() string {
	return r.c.Nationality
}

func (r *celebrityResolver) ReqPhotoPath() string {
	return r.c.ReqPhotoPath
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

func openAndInitDatabase(connString string) (*gorm.DB, error) {
	masterConnString := strings.Replace(connString, "database="+databaseName(), "database=master", 1)
	masterDB, err := sql.Open("sqlserver", masterConnString)
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

	gormDB, err := gorm.Open(sqlserver.Open(connString), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Warn,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		_ = sqlDB.Close()
		return nil, err
	}

	return gormDB, nil
}

func connectionCandidates(dbName string) []string {
	if raw := strings.TrimSpace(os.Getenv("GO10_CONNECTION_STRING")); raw != "" {
		return []string{raw}
	}

	server := envOrDefault("GO10_SQLSERVER", defaultSQLHost)
	port := envOrDefault("GO10_SQLPORT", defaultSQLPort)
	base := []string{
		buildConnectionString(server, port, dbName),
		buildConnectionString(`localhost`, `65000`, dbName),
		buildConnectionString(`MARINAPC\MARINAPC`, ``, dbName),
		buildConnectionString(`localhost\MARINAPC`, ``, dbName),
		buildConnectionString(`localhost`, `64654`, dbName),
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
	return envOrDefault("GO10_DATABASE", defaultDatabase)
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
