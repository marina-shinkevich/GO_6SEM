package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/graph-gophers/graphql-go"
)


const schemaString = `
	schema {
		query: Query
	}
	type Query {
		celebrity(id: ID!): Celebrity
	}
	type Celebrity {
		id: ID!
		firstname: String!
		surname: String!
	}
`

type Celebrity struct {
	ID        graphql.ID
	Firstname string
	Surname   string
}

// === 3. РЕЗОЛВЕРЫ (Resolvers) ===
// Резолвер — это структура, методы которой физически добывают данные для GraphQL [2]
type Resolver struct{}

type CelebrityResolver struct {
	c *Celebrity
}

// Обработчик корневого запроса "celebrity(id: ID!)"
func (r *Resolver) Celebrity(ctx context.Context, args struct{ ID graphql.ID }) *CelebrityResolver {
	// В реальном проекте здесь был бы поиск в БД.
	// В учебном примере просто возвращаем объект с переданным ID.
	return &CelebrityResolver{&Celebrity{ID: args.ID}}
}

// Резолверы конкретных полей (возвращаем заглушки "stub", как на слайде лекции) [6]
func (c *CelebrityResolver) ID() graphql.ID    { return c.c.ID }
func (c *CelebrityResolver) Firstname() string { return "firstname-stub" }
func (c *CelebrityResolver) Surname() string   { return "surname-stub" }

// === 4. HTTP-ОБРАБОТЧИК ДЛЯ POSTMAN ===
func graphqlHandler(schema *graphql.Schema) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Postman отправляет GraphQL-запрос в виде JSON-объекта
		var params struct {
			Query         string                 `json:"query"`
			OperationName string                 `json:"operationName"`
			Variables     map[string]interface{} `json:"variables"`
		}

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Выполняем распарсенный запрос через движок GraphQL [3]
		response := schema.Exec(r.Context(), params.Query, params.OperationName, params.Variables)

		// Отправляем JSON-ответ обратно в Postman
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(response)
	}
}

func main() {
	// Загружаем и парсим схему, связывая её с нашим корневым Резолвером [3]
	schema := graphql.MustParseSchema(schemaString, &Resolver{})

	// Регистрируем единственный URL (endpoint) для всех запросов
	http.HandleFunc("/graphql", graphqlHandler(schema))

	log.Println("GraphQL server is running on http://localhost:3000/graphql")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
