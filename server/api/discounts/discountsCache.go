package discounts

import (
	"encoding/json"

	"github.com/aiescola/go_shop_api/util"
	"github.com/go-redis/redis/v8"
)

type discountsRedisCache struct {
	util.RedisCache
}

type DiscountsCache interface {
	GetAll() ([]Discount, error)
	PutAll(discounts []Discount) error
}

func NewDiscountsRedisCache(client *redis.Client) DiscountsCache {
	return discountsRedisCache{
		util.NewRedisCache(client, "Discounts"),
	}
}

func (r discountsRedisCache) GetAll() ([]Discount, error) {
	data, err := r.Client.Get(r.Ctx, r.Key).Result()
	if err != nil {
		return nil, err
	}

	var discounts []Discount
	if err := json.Unmarshal([]byte(data), &discounts); err != nil {
		return nil, err
	}

	return discounts, nil
}

func (r discountsRedisCache) PutAll(discounts []Discount) error {
	jsonDiscounts, err := json.Marshal(discounts)
	if err != nil {
		return err
	}

	if err := r.Client.Set(r.Ctx, r.Key, jsonDiscounts, r.Ttl); err != nil {
		return err.Err()
	}
	return nil
}
