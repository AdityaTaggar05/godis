package store

import "sync"

var DB sync.Map

func Load(key string) (Data, bool) {
	value, found := DB.Load(key)

	if !found {
		return Data{}, false
	}

	data := value.(Data)

	if data.isExpired() {
		Delete(key)
		return Data{}, false
	}

	return data, true
}

func Store(key string, d Data) {
	DB.Store(key, d)
}

func LoadOrStore(key string, d Data) (Data, bool) {
	value, found := Load(key)

	if found {
		return value, true
	}

	Store(key, d)
	return d, false
}

func Delete(key string) {
	DB.Delete(key)
}
