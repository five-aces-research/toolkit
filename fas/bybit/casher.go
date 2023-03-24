package bybit

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"
)

type Cache struct {
	data map[string]CacheEntry
}

type CacheEntry struct {
	timestamp time.Time
	value     interface{}
}

func (c *Cache) Get(key string, validFor int) (interface{}, bool) {
	if entry, ok := c.data[key]; ok {
		if time.Since(entry.timestamp) < time.Duration(validFor)*time.Second {
			return entry.value, true
		}
	}
	return nil, false
}

func (c *Cache) Set(key string, value interface{}) {
	c.data[key] = CacheEntry{
		timestamp: time.Now(),
		value:     value,
	}
}

func (b *Public) Check(validFor int, result interface{}, apiFunc interface{}, parameters ...interface{}) error {
	funcType := reflect.TypeOf(apiFunc)
	if funcType.Kind() != reflect.Func {
		return fmt.Errorf("apiFunc must be a function")
	}

	hash := createHash(apiFunc, parameters)
	if value, ok := b.cache.Get(hash, validFor); ok {
		fmt.Println("wow")
		reflect.ValueOf(result).Elem().Set(reflect.ValueOf(value).Elem())
	} else {
		// Call the API function
		fmt.Println("here first")
		res := reflect.ValueOf(apiFunc).Call(convertToValue(parameters))
		if err, ok := res[1].Interface().(error); ok {
			// use err variable here
			return err
		}
		// Save data to the cache
		b.cache.Set(hash, res[0].Interface())
		// Assign the result

		reflect.ValueOf(result).Elem().Set(res[0].Elem())
	}

	return nil
}

func createHash(apiFunc interface{}, parameters []interface{}) string {
	funcName := runtime.FuncForPC(reflect.ValueOf(apiFunc).Pointer()).Name()
	funcName = strings.Replace(funcName, ".", "_", -1)
	hash := funcName
	for _, param := range parameters {
		hash += fmt.Sprintf("_%v", param)
	}
	return hash
}

func convertToValue(params []interface{}) []reflect.Value {
	values := make([]reflect.Value, len(params))
	for i, param := range params {
		values[i] = reflect.ValueOf(param)
	}
	return values
}
