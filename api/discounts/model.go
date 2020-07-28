package discounts

type Discount struct {
	Code         string   `json:"code"`
	Type         string   `json:"type"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	ProductCodes []string `json:"productCodes"`
	// For products
	Price float64 `json:"price,omitempty"`
	// For bulk
	Pct       int `json:"pct,omitempty"`
	MinAmount int `json:"minAmount,omitempty"`
	// For promotions
	Buy int `json:"buy,omitempty"`
	Pay int `json:"pay,omitempty"`
}
