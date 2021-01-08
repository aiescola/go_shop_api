package discounts

import (
	"encoding/json"
	"net/http"

	"github.com/aiescola/go_shop_api/util"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type DiscountController struct {
	dataSource     DiscountDataSource
	discountsCache DiscountsCache
	logger         *log.Entry
}

func NewController(ds DiscountDataSource, redisDataSource DiscountsCache, logger *log.Logger) DiscountController {
	return DiscountController{
		ds,
		redisDataSource,
		logger.WithFields(log.Fields{
			"file": "DiscountController",
		}),
	}
}

func (dc DiscountController) GetDiscountOld(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	var discounts []Discount

	if err := dc.getDiscountsInternal(&discounts); err != nil {
		dc.logger.Error(err.Error())
		util.EncodeError(response, http.StatusInternalServerError, err.Error())
		return
	}

	dc.logger.Info("Discounts retrieved: ", discounts)

	json.NewEncoder(response).Encode(discounts[0])
}

func (dc DiscountController) GetDiscounts(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")

	var discounts []Discount

	if err := dc.getDiscountsInternal(&discounts); err != nil {
		dc.logger.Error(err.Error())
		util.EncodeError(response, http.StatusInternalServerError, err.Error())
		return
	}

	dc.logger.Info("Discounts retrieved: ", discounts)

	json.NewEncoder(response).Encode(struct {
		Discounts []Discount `json:"discounts"`
	}{discounts})
}

func (dc DiscountController) getDiscountsInternal(discounts *[]Discount) error {
	var err error
	if *discounts, err = dc.discountsCache.GetAll(); err != nil {
		*discounts, err = dc.dataSource.GetDiscounts()
		dc.logger.Info("Retrieving discounts from mongo")

		if err != nil {
			return err
		}

		if err := dc.discountsCache.PutAll(*discounts); err != nil {
			dc.logger.Error(err.Error())
		}
	}

	return nil
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

	// check that the discount doesn't already exist on DB
	if dbDiscount, _ := dc.dataSource.GetDiscount(discount.Code); dbDiscount != nil {
		dc.logger.Error("Discount already exists")
		util.EncodeError(response, http.StatusBadRequest, "Discount already exists")
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
