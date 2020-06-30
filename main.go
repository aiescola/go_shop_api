package main

import (
	"apirest/products"
	"net/http"

	"github.com/gorilla/mux"
)

type ApiService interface {
	AddRoutes(router *mux.Router)
}

var productService ApiService

func main() {
	dataSource := products.MakeLocalProductDataSource()
	productService = products.New(dataSource)

	router := mux.NewRouter()
	productService.AddRoutes(router)

	http.ListenAndServe(":8080", router)
}
