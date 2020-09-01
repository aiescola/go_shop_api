package products

import (
	"encoding/json"

	"github.com/aiescola/go_shop_api/util"
	"github.com/go-redis/redis/v8"
)

type productsRedisCache struct {
	util.RedisCache
}

type ProductsCache interface {
	GetOne(code string) (*Product, error)
	GetAll() ([]Product, error)
	PutOne(product Product) error
	PutAll(products []Product) error
}

func NewProductsRedisCache(client *redis.Client) ProductsCache {
	return productsRedisCache{
		util.NewRedisCache(client, "Products"),
	}
}

func (r productsRedisCache) GetOne(code string) (*Product, error) {
	data, err := r.Client.Get(r.Ctx, code).Result()
	if err != nil {
		return nil, err
	}

	var product Product
	if err := json.Unmarshal([]byte(data), &product); err != nil {
		return nil, err
	}
	return &product, nil
}

func (r productsRedisCache) GetAll() ([]Product, error) {
	data, err := r.Client.Get(r.Ctx, r.Key).Result()
	if err != nil {
		return nil, err
	}

	var products []Product
	if err := json.Unmarshal([]byte(data), &products); err != nil {
		return nil, err
	}

	return products, nil
}

func (r productsRedisCache) PutOne(product Product) error {
	jsonProduct, err := json.Marshal(product)
	if err != nil {
		return err
	}

	if err := r.Client.Set(r.Ctx, product.Code, jsonProduct, r.Ttl); err != nil {
		return err.Err()
	}
	return nil
}

func (r productsRedisCache) PutAll(products []Product) error {
	jsonProducts, err := json.Marshal(products)
	if err != nil {
		return err
	}

	if err := r.Client.Set(r.Ctx, r.Key, jsonProducts, r.Ttl); err != nil {
		return err.Err()
	}
	return nil
}
