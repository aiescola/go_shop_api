package main

import (
	"apirest/products"
	"apirest/util"
	"net/http"

	"github.com/gorilla/mux"
)

type ApiService interface {
	AddRoutes(router *mux.Router)
}

var productService ApiService

func main() {
	port := util.GetEnv("PORT", "8080")

	dataSource := products.MakeLocalProductDataSource()
	productService = products.New(dataSource)

	router := mux.NewRouter()
	productService.AddRoutes(router)

	http.ListenAndServe(":"+port, router)
}
