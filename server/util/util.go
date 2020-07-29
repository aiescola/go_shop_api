package util

import (
	"os"
)

func GetEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		if len(value) > 0 {
			return value
		}
	}
	return defaultValue
}
