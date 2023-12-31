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
	router.HandleFunc("/user", handlers.CreateUserHandler).Methods(http.MethodPost)
	router.HandleFunc("/cart/{mobile_number}", handlers.AddItemToCart).Methods(http.MethodPost)
	router.HandleFunc("/cart/remove/{mobile_number}", handlers.RemoveItemFromCart).Methods(http.MethodPost)
	router.HandleFunc("/cart/{mobile_number}", handlers.ViewCart).Methods(http.MethodGet)
	router.HandleFunc("/order/create/{mobile_number}", handlers.CreateOrder).Methods(http.MethodPost)
	router.HandleFunc("/order/utm", handlers.ViewOrdersByUTM).Methods(http.MethodGet)
	return router
}
