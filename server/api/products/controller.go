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
	dataSource    ProductDataSource
	productsCache ProductsCache
	logger        *log.Entry
}

//NewController using the dataSource and logger passed as parameters
func NewController(ds ProductDataSource, productsCache ProductsCache, logger *log.Logger) ProductController {
	return ProductController{
		ds,
		productsCache,
		logger.WithFields(log.Fields{
			"file": "ProductController",
		}),
	}
}

func (p ProductController) GetProducts(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	var products []Product
	var err error
	if products, err = p.productsCache.GetAll(); err != nil {
		products, err = p.dataSource.GetProducts()
		p.logger.Info("Retrieving products from mongo")

		if err != nil {
			p.logger.Error(err.Error())
			util.EncodeError(response, http.StatusInternalServerError, err.Error())
			return
		}

		if err := p.productsCache.PutAll(products); err != nil {
			p.logger.Error(err.Error())
		}
	}

	p.logger.Info("Products retrieved: ", products)

	json.NewEncoder(response).Encode(struct {
		Products []Product `json:"products"`
	}{products})
}

func (p ProductController) GetProduct(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	code := mux.Vars(request)["code"]

	var product *Product
	var err error

	if product, err = p.productsCache.GetOne(code); err != nil {
		product, err = p.dataSource.GetProduct(code)
		p.logger.Info("Retrieving product from mongo")

		if err != nil {
			p.logger.Error(err)
			util.EncodeError(response, http.StatusInternalServerError, err.Error())
			return
		}

		if err := p.productsCache.PutOne(*product); err != nil {
			p.logger.Error(err.Error())
		}
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
