package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"

	"github.com/gorilla/mux"
)

type RPCRequest struct {
	JSONRPC string           `json:"jsonrpc"`
	Method  string           `json:"method"`
	Params  json.RawMessage  `json:"params,omitempty"`
	ID      *json.RawMessage `json:"id,omitempty"`
}

type RPCResponse struct {
	JSONRPC string           `json:"jsonrpc"`
	Result  interface{}      `json:"result,omitempty"`
	Error   *RPCError        `json:"error,omitempty"`
	ID      *json.RawMessage `json:"id"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

const (
	ErrParseError     = -32700
	ErrInvalidRequest = -32600
	ErrMethodNotFound = -32601
	ErrInvalidParams  = -32602
	ErrInternalError  = -32603
)

var precision = 2

func rpcError(code int, msg string) *RPCError {
	return &RPCError{Code: code, Message: msg}
}

func makeResponse(id *json.RawMessage, result interface{}, err *RPCError) *RPCResponse {
	return &RPCResponse{
		JSONRPC: "2.0",
		Result:  result,
		Error:   err,
		ID:      id,
	}
}

func parseXY(params json.RawMessage) (float64, float64, error) {

	var arr []float64
	if err := json.Unmarshal(params, &arr); err == nil {
		if len(arr) != 2 {
			return 0, 0, fmt.Errorf("expected 2 params, got %d", len(arr))
		}
		return arr[0], arr[1], nil
	}

	var obj struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	}
	if err := json.Unmarshal(params, &obj); err == nil {
		return obj.X, obj.Y, nil
	}

	return 0, 0, fmt.Errorf("invalid params format")
}

func parseN(params json.RawMessage) (int, error) {
	var obj struct {
		N int `json:"N"`
	}
	if err := json.Unmarshal(params, &obj); err != nil {
		return 0, fmt.Errorf("invalid params for pre: %v", err)
	}
	return obj.N, nil
}

func formatFloat(v float64) interface{} {
	factor := math.Pow(10, float64(precision))
	rounded := math.Round(v*factor) / factor
	return rounded
}

func handleSingle(req RPCRequest) *RPCResponse {
	isNotification := req.ID == nil

	if req.JSONRPC != "2.0" {
		if isNotification {
			return nil
		}
		return makeResponse(req.ID, nil, rpcError(ErrInvalidRequest, "invalid jsonrpc version, expected \"2.0\""))
	}

	log.Printf("[RPC] method=%q id=%v notification=%v", req.Method, req.ID, isNotification)

	switch req.Method {
	case "sum", "sub", "mul", "div":
		x, y, err := parseXY(req.Params)
		if err != nil {
			if isNotification {
				return nil
			}
			return makeResponse(req.ID, nil, rpcError(ErrInvalidParams, err.Error()))
		}

		var result float64
		switch req.Method {
		case "sum":
			result = x + y
		case "sub":
			result = x - y
		case "mul":
			result = x * y
		case "div":
			if y == 0 {
				if isNotification {
					return nil
				}
				return makeResponse(req.ID, nil, rpcError(ErrInvalidParams, "division by zero"))
			}
			result = x / y
		}

		log.Printf("[RPC] %s(%v, %v) = %v (precision=%d)", req.Method, x, y, result, precision)

		if isNotification {
			return nil
		}
		return makeResponse(req.ID, formatFloat(result), nil)

	case "pre":
		n, err := parseN(req.Params)
		if err != nil {
			return nil
		}
		precision = n
		log.Printf("[RPC] pre: precision set to %d", precision)
		return nil

	default:
		if isNotification {
			return nil
		}
		return makeResponse(req.ID, nil, rpcError(ErrMethodNotFound, fmt.Sprintf("method not found: %s", req.Method)))
	}
}

func rpcHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var raw json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
		log.Printf("[RPC] parse error: %v", err)
		resp := makeResponse(nil, nil, rpcError(ErrParseError, "parse error"))
		json.NewEncoder(w).Encode(resp)
		return
	}

	trimmed := raw
	isBatch := len(trimmed) > 0 && trimmed[0] == '['

	if isBatch {
		var reqs []RPCRequest
		if err := json.Unmarshal(raw, &reqs); err != nil {
			resp := makeResponse(nil, nil, rpcError(ErrParseError, "batch parse error"))
			json.NewEncoder(w).Encode(resp)
			return
		}

		if len(reqs) == 0 {
			resp := makeResponse(nil, nil, rpcError(ErrInvalidRequest, "empty batch"))
			json.NewEncoder(w).Encode(resp)
			return
		}

		log.Printf("[RPC] batch request, %d items", len(reqs))

		var responses []*RPCResponse
		for _, req := range reqs {
			resp := handleSingle(req)
			if resp != nil {
				responses = append(responses, resp)
			}
		}

		if len(responses) == 0 {

			w.WriteHeader(http.StatusNoContent)
			return
		}

		json.NewEncoder(w).Encode(responses)
		return
	}

	var req RPCRequest
	if err := json.Unmarshal(raw, &req); err != nil {
		resp := makeResponse(nil, nil, rpcError(ErrInvalidRequest, "invalid request"))
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := handleSingle(req)
	if resp == nil {

		w.WriteHeader(http.StatusNoContent)
		return
	}

	json.NewEncoder(w).Encode(resp)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/rpc", rpcHandler).Methods(http.MethodPost)

	log.Println("JSON-RPC 2.0 server starting on :3000")
	if err := http.ListenAndServe(":3000", r); err != nil {
		log.Fatal(err)
	}
}
