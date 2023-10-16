package main

import(
	"shopify-apis/server"
	"github.com/urfave/negroni"
	"shopify-apis/configs"
)
func main() {
	configs.Load()
	router := server.InitRouter()

	server := negroni.Classic()
	server.UseHandler(router)
	server.Run(":" + configs.AppPort())
}