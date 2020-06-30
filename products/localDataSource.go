package products

type localDataSource struct {
	products []Product
}

func (d *localDataSource) GetProducts() ([]Product, error) {
	return d.products, nil
	//	return nil, errors.New("Failure getting product")
}

func (d *localDataSource) AddProduct(product Product) {
	d.products = append(d.products, product)
}
