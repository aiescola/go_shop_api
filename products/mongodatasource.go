package products

/*
type fakeProductDataSource struct {
	products []Product
}

func (d *fakeProductDataSource) GetProducts() []Product {
	return d.products
}

func (d *fakeProductDataSource) AddProduct(product Product) {
	d.products = append(d.products, product)
}
*/

type mongoProductDataSource struct{}

// func (d *mongoProductDataSource) GetProducts() ([]Product, error) {
// 	//return d.products
// }

// func (d *mongoProductDataSource) AddProduct(product Product) {
// 	//d.products = append(d.products, product)
// }
