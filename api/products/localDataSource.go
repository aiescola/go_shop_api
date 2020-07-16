package products

import "errors"

type localDataSource struct {
	products []Product
}

func (d *localDataSource) GetProducts() ([]Product, error) {
	return d.products, nil
	//	return nil, errors.New("Failure getting product")
}

func (d *localDataSource) GetProduct(code string) (*Product, error) {
	for _, product := range d.products {
		if product.Code == code {
			return &product, nil
		}
	}
	return nil, errors.New("Product not found")
}

func (d *localDataSource) AddProduct(product Product) {
	d.products = append(d.products, product)
}
