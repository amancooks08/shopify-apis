package server

import (
	"net/http"
	"github.com/gorilla/mux"
	"shopify-apis/service"
)

func InitRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/ping", service.PingHandler).Methods(http.MethodGet)
	return router
}