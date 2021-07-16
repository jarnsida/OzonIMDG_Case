package server

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type memoryDB struct {
	items map[string]string
	mu    sync.RWMutex
}

func newDB() memoryDB {
	f, err := os.Open("./store/db.json")
	if err != nil {
		return memoryDB{items: map[string]string{}}
	}
	items := map[string]string{}
	if err := json.NewDecoder(f).Decode(&items); err != nil {
		fmt.Println("не удалось декодировать файл", err.Error())
		return memoryDB{items: map[string]string{}}
	}
	return memoryDB{items: items}
}

func (m *memoryDB) set(key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items[key] = value
}

func (m *memoryDB) get(key string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, found := m.items[key]
	return value, found
}

func (m *memoryDB) delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.items, key)
}

func (m *memoryDB) save() {
	f, err := os.Create("./store/db.json")
	if err != nil {
		fmt.Println("не удалось создать файл", err.Error())
	}
	if err := json.NewEncoder(f).Encode(m.items); err != nil {
		fmt.Println("не удалось закодировать", err.Error())
	}
	fmt.Println("успешное сохранение db в файл")
}
