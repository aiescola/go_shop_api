package products

type ProductDataSource interface {
	GetProducts() ([]Product, error)
	AddProduct(product Product)
}

func MakeLocalProductDataSource() *localDataSource {
	return &localDataSource{[]Product{{"MUG", "Mug", 5.4}, {"TSHIRT", "T-Shirt", 12.5}}}
}
