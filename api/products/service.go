package products

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

//ProductService
type ProductService struct {
	dataSource ProductDataSource
}

type ErrorResponse struct {
	Status int    `json:"status"`
	Err    string `json:"error"`
}

//NewService ProductService using the dataSource passed as parameter
func NewService(ds ProductDataSource) *ProductService {
	return &ProductService{
		dataSource: ds,
	}
}

func (p *ProductService) AddRoutes(router *mux.Router) {
	router.HandleFunc("/api/products", p.getProducts).Methods("GET")
	router.HandleFunc("/api/products/{Code}", p.getProduct).Methods("GET")
	router.HandleFunc("/api/products", p.addProduct).Methods("POST")
}

func (p *ProductService) getProducts(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	products, err := p.dataSource.GetProducts()

	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(ErrorResponse{http.StatusInternalServerError, err.Error()})
		return
	}

	json.NewEncoder(response).Encode(struct {
		Products []Product `json:"products"`
	}{products})
}

func (p *ProductService) getProduct(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)

	product, err := p.dataSource.GetOne(params["Code"])
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode(ErrorResponse{http.StatusInternalServerError, err.Error()})
		return
	}

	json.NewEncoder(response).Encode(product)
}

func (p *ProductService) addProduct(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var product Product

	err := json.NewDecoder(request.Body).Decode(&product)

	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(ErrorResponse{http.StatusBadRequest, err.Error()})
		return
	}

	if product.Code == "" || product.Name == "" {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(ErrorResponse{http.StatusBadRequest, "Invalid body format"})
		return
	}

	err = p.dataSource.AddProduct(product)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(ErrorResponse{http.StatusInternalServerError, err.Error()})
		return
	}

	json.NewEncoder(response).Encode("Item added")
}
