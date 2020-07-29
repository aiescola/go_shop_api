package discounts

import (
	"encoding/json"
	"net/http"

	"github.com/aiescola/go_shop_api/util"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type DiscountController struct {
	dataSource DiscountDataSource
	logger     *log.Entry
}

func NewController(ds DiscountDataSource, logger *log.Logger) DiscountController {
	return DiscountController{
		ds,
		logger.WithFields(log.Fields{
			"file": "DiscountController",
		}),
	}
}

func (dc DiscountController) GetDiscounts(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	discounts, err := dc.dataSource.GetDiscounts()

	if err != nil {
		dc.logger.Error(err.Error())
		util.EncodeError(response, http.StatusInternalServerError, err.Error())
		return
	}

	json.NewEncoder(response).Encode(struct {
		Discounts []Discount `json:"discounts"`
	}{discounts})
}

func (dc DiscountController) GetDiscount(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	params := mux.Vars(request)

	discount, err := dc.dataSource.GetDiscount(params["code"])

	if err != nil {
		dc.logger.Error(err.Error())
		util.EncodeError(response, http.StatusInternalServerError, err.Error())
		return
	}

	json.NewEncoder(response).Encode(discount)
}

func (dc DiscountController) AddDiscount(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	var discount Discount

	if err := json.NewDecoder(request.Body).Decode(&discount); err != nil {
		dc.logger.Error(err)
		util.EncodeError(response, http.StatusBadRequest, err.Error())
		return
	}

	if discount.Code == "" || discount.Name == "" {
		dc.logger.Error("Invalid body format")
		util.EncodeError(response, http.StatusBadRequest, "Invalid body format")
		return
	}

	if err := dc.dataSource.AddDiscount(discount); err != nil {
		dc.logger.Error(err)
		util.EncodeError(response, http.StatusBadRequest, err.Error())
		return
	}

	json.NewEncoder(response).Encode(struct {
		Discount Discount `json:"Discount"`
	}{discount})

}
