package main

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"
)



var (
	
	methods = map[string]func(RPCRequest) RPCResponse{}
)


type RPCRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params json.RawMessage `json:"params"`
	ID     interface{}     `json:"id"` 
}

// RPCResponse — структура ответа JSON-RPC 2.0
type RPCResponse struct {
	Jsonrpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"` // Опускается в JSON, если пусто
	Error   interface{} `json:"error,omitempty"`  // Опускается в JSON, если пусто
	ID      interface{} `json:"id"`
}

// Структуры для конкретного метода (например, сложения)
type SumArg struct{ X, Y float64 }
type SumResult struct{ Value float64 }

// === 2. ЯДРО RPC-СЕРВЕРА ===

// register — функция для регистрации новых методов в словаре сервера
func register(name string, fn func(RPCRequest) RPCResponse) {
	methods[name] = fn
}

// callMethod — диспетчер, который ищет метод по имени и вызывает его
func callMethod(r RPCRequest) RPCResponse {
	if method := methods[r.Method]; method != nil {
		return method(r) // Вызываем найденную функцию
	}
	// Если клиент запросил несуществующий метод, возвращаем ошибку
	return RPCResponse{
		Jsonrpc: "2.0",
		Error:   "method not found",
		ID:      r.ID,
	}
}

// rpcHandler — единый HTTP-обработчик для ВСЕХ запросов к серверу
func rpcHandler(w http.ResponseWriter, r *http.Request) {
	var raw json.RawMessage

	// Читаем тело HTTP-запроса
	if err := json.NewDecoder(r.Body).Decode(&raw); err == nil {

		// Проверка на пакетный (batch) запрос: начинается ли JSON с символа массива '['
		if len(raw) > 0 && raw[0] == '[' {
			// Логика пакетной обработки (упрощена для базового примера)
			// var reqs []RPCRequest
			// json.Unmarshal(raw, &reqs) ...
		} else { // Одиночный запрос
			var rq RPCRequest
			json.Unmarshal(raw, &rq) // Декодируем основную структуру запроса

			if rq.ID == nil { // notification (уведомление - запрос без ответа)
				callMethod(rq)
			} else { // обычный запрос, на который клиент ждет ответ
				response := callMethod(rq)
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(response)
			}
		}
	} else {
		// Ошибка разбора самого JSON-конверта
		json.NewEncoder(w).Encode(RPCResponse{Jsonrpc: "2.0", Error: "json error"})
	}
}

// === 3. ТОЧКА ВХОДА ===

func main() {
	// 1. Регистрируем бизнес-логику (метод "sum")
	register("sum", func(r RPCRequest) RPCResponse {
		log.Println("Вызван метод sum")
		var p SumArg

		// В лекциях здесь используется сложная кастомная функция extractParams,
		// но для базовой механики достаточно стандартного парсера:
		json.Unmarshal(r.Params, &p)
		return RPCResponse{
			Jsonrpc: "2.0",
			Result:  SumResult{Value: p.X + p.Y},
			ID:      r.ID,
		}
	})

	router := mux.NewRouter()

	router.HandleFunc("/rpc", rpcHandler).Methods("POST")

	log.Println("JSON-RPC Server started on port 3000")
	http.ListenAndServe(":3000", router)
}
