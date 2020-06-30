package products

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

//ProductApi
type ProductApi struct {
	dataSource ProductDataSource
}

//New ProductApi
func New(ds ProductDataSource) *ProductApi {
	return &ProductApi{
		dataSource: ds,
	}
}

func (p *ProductApi) AddRoutes(router *mux.Router) {
	router.HandleFunc("/api/products", p.getProducts).Methods("GET")
	router.HandleFunc("/api/products/{Code}", p.getProduct).Methods("GET")
	router.HandleFunc("/api/products", p.addProduct).Methods("POST")
}

func (p *ProductApi) getProducts(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	products, err := p.dataSource.GetProducts()

	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(response).Encode(struct {
		Products []Product `json:"products"`
	}{products})
}

func (p *ProductApi) getProduct(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)

	products, err := p.dataSource.GetProducts()
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, product := range products {
		if product.Code == params["Code"] {
			response.Header().Set("Content-Type", "application/json")
			json.NewEncoder(response).Encode(product)
		}
	}

}

func (p *ProductApi) addProduct(w http.ResponseWriter, r *http.Request) {
	var product Product

	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		//TODO: return error
		return
	}

	p.dataSource.AddProduct(product)
}
