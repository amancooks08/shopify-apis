package server

import (
	"net/http"
	"shopify-apis/handlers"

	"github.com/gorilla/mux"
)

func InitRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/ping", handlers.PingHandler).Methods(http.MethodGet)
	router.HandleFunc("/products", handlers.GetProductsHandler).Methods(http.MethodGet)
	return router
}
