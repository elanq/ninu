package ninu

import (
	"errors"
	"time"

	"github.com/allegro/bigcache/v2"
)

const timeRepresentation = 15

var (
	CacheExpired = errors.New("Cache expired")
	Redis        Cache
)

func InitRedis() {
	Redis = NewRedisCache()
}

type Cache interface {
	Set(key string, val []byte, t ...time.Duration) error
	Get(key string) ([]byte, error)
}

func NewMemoryCache() Cache {
	conf := bigcache.Config{
		Shards:           2,
		LifeWindow:       10 * time.Minute,
		HardMaxCacheSize: 50,
	}
	cacheClient, err := bigcache.NewBigCache(conf)
	if err != nil {
		panic(err)
	}

	return &MemoryCache{
		client: cacheClient,
	}
}

type MemoryCache struct {
	client *bigcache.BigCache
}

func (m *MemoryCache) Set(key string, val []byte, t ...time.Duration) error {
	var expireIn time.Time
	if len(t) == 1 {
		expireIn = time.Now().Add(t[0])
	}
	exprBinary, err := expireIn.MarshalBinary()
	if err != nil {
		return err
	}
	v := append(exprBinary, val...)
	return m.client.Set(key, v)
}

func (m *MemoryCache) Get(key string) ([]byte, error) {
	var expireTime time.Time
	value, err := m.client.Get(key)
	if err != nil {
		return nil, err
	}

	err = expireTime.UnmarshalBinary(value[:timeRepresentation])
	if err != nil {
		return nil, err
	}
	//doesn't have expiration time
	if (expireTime == time.Time{}) {
		return value[timeRepresentation:], nil
	}

	if time.Now().After(expireTime) {
		return nil, CacheExpired
	}
	return value[timeRepresentation:], nil
}
