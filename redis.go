package ninu

import (
	"errors"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

func NewRedisCache() KeyValueCache {
	client, err := redis.DialURL(os.Getenv("REDIS_URL"))
	if err != nil {
		panic(err)
	}
	return &RedisCache{
		client: client,
	}
}

type RedisCache struct {
	client redis.Conn
}

func (r *RedisCache) Set(key string, val []byte, t ...time.Duration) error {
	args := []interface{}{
		key, val,
	}
	if len(t) == 1 {
		duration := int64(t[0].Seconds())
		args = append(args, "EX", duration)
	}
	_, err := r.client.Do("SET", args...)
	return err
}

func (r *RedisCache) Get(key string) ([]byte, error) {
	reply, err := r.client.Do("GET", key)
	if err != nil {
		return nil, err
	}

	if reply == nil {
		return nil, nil
	}
	if val, ok := reply.([]byte); ok {
		return val, nil
	}
	return nil, errors.New("Error while parsing value")
}
