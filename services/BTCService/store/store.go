package store

import (
	"fmt"
	"os"
	"encoding/gob"

	"xox_tokens/services/BTCService/configs"

	"github.com/patrickmn/go-cache"
)

var store *cache.Cache

// Initiate store.
func InitiateStore(storePath string) {
	fmt.Println("Initiate store...")

	items, err := getFromFile(storePath)
	if err != nil {
		fmt.Println("failed to rise store from file", err)
	}

	if len(items) > 0 {
		fmt.Println("cache items:", items)
		store = cache.NewFrom(cache.NoExpiration, cache.NoExpiration, items)
	} else {
		store = cache.New(cache.NoExpiration, cache.NoExpiration)
	}

	return
}

// Setter for cache.
// Sets value without expiration.
func Set(key string, value interface{}) {
	//fmt.Println("Store - set...", key, value)
	store.Set(key, value, cache.NoExpiration)
}

// Getter for cache.
func Get(key string) (value interface{}, found bool) {
	value, found = store.Get(key)
	//fmt.Println("Store - get", key, value)
	return
}

// Save saves store items to file.
func Save() (err error) {
	fmt.Println("Saving cache...")
	conf := configs.GetConfigs()
	items := store.Items()

	//fmt.Println("cache items:", items)

	f, err := os.OpenFile(conf.Common.StorePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println("ERROR OPEN FILE:", err.Error())
		return
	}
	defer f.Close()

	encoder := gob.NewEncoder(f)
	err = encoder.Encode(items)

	return
}

func getFromFile(storePath string) (map[string]cache.Item, error) {
	var items map[string]cache.Item

	fmt.Println("Getting cache from file:", storePath)
	f, err := os.OpenFile(storePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil { return items, err }
	defer f.Close()

	decoder := gob.NewDecoder(f)
	err = decoder.Decode(&items)

	return items, err
}
