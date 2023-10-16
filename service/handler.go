package service

import(
	"net/http"
	"encoding/json"
	"shopify-apis/domain"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	response := domain.PingResponse{Message: "pong"}
	json.NewEncoder(w).Encode(response)
}