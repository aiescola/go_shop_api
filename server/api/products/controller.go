package products

import (
	"encoding/json"
	"net/http"

	"github.com/aiescola/go_shop_api/util"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

//ProductController
type ProductController struct {
	dataSource ProductDataSource
	logger     *log.Entry
}

//NewController using the dataSource and logger passed as parameters
func NewController(ds ProductDataSource, logger *log.Logger) ProductController {
	return ProductController{
		ds,
		logger.WithFields(log.Fields{
			"file": "ProductController",
		}),
	}
}

func (p ProductController) GetProducts(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	products, err := p.dataSource.GetProducts()

	if err != nil {
		p.logger.Error(err.Error())
		util.EncodeError(response, http.StatusInternalServerError, err.Error())
		return
	}

	p.logger.Info("Products retrieved: ", products)

	json.NewEncoder(response).Encode(struct {
		Products []Product `json:"products"`
	}{products})
}

func (p ProductController) GetProduct(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	params := mux.Vars(request)

	product, err := p.dataSource.GetProduct(params["code"])
	if err != nil {
		p.logger.Error(err)
		util.EncodeError(response, http.StatusInternalServerError, err.Error())
		return
	}

	json.NewEncoder(response).Encode(product)

}

func (p ProductController) AddProduct(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var product Product

	if err := json.NewDecoder(request.Body).Decode(&product); err != nil {
		p.logger.Error(err)
		util.EncodeError(response, http.StatusBadRequest, err.Error())
		return
	}

	if product.Code == "" || product.Name == "" {
		p.logger.Error("Invalid body format")
		util.EncodeError(response, http.StatusBadRequest, "Invalid body format")
		return
	}

	// check that the product doesn't already exist on DB
	if dbProduct, _ := p.dataSource.GetProduct(product.Code); dbProduct != nil {
		p.logger.Error("Product already exists")
		util.EncodeError(response, http.StatusBadRequest, "Product already exists")
		return
	}

	if err := p.dataSource.AddProduct(product); err != nil {
		p.logger.Error(err)
		util.EncodeError(response, http.StatusBadRequest, err.Error())
		return
	}

	json.NewEncoder(response).Encode("Item added")
}
