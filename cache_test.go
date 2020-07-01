package ninu_test

import (
	"testing"
	"time"

	"github.com/elanq/ninu"
)

func TestMemoryCache(t *testing.T) {
	cache := ninu.NewMemoryCache()
	cases := []struct {
		key            string
		value          string
		expirationTime time.Duration
	}{
		{
			key:   "something-unique",
			value: "some value",
		},
		{
			key:            "something-unique",
			value:          "some value",
			expirationTime: 1 * time.Millisecond,
		},
		{
			key:            "something-unique",
			value:          "some value",
			expirationTime: 1 * time.Nanosecond,
		},
	}

	for _, c := range cases {
		var err error
		if c.expirationTime == 0 {
			err = cache.Set(c.key, []byte(c.value))
		} else {
			err = cache.Set(c.key, []byte(c.value), c.expirationTime)
		}

		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		//test expired cache
		if c.expirationTime == (1 * time.Nanosecond) {
			_, err := cache.Get(c.key)
			if err != ninu.CacheExpired {
				t.Error(err)
				t.FailNow()
			}
			return
		}

		val, err := cache.Get(c.key)
		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		if string(val) != c.value {
			t.Error("unexpected cache value")
			t.Errorf("expected %v, actual %v \n", c.value, string(val))
			t.FailNow()
		}
	}
}
