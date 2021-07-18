package server

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

type memoryDB struct {
	mu    sync.RWMutex
	items map[string]string
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

func (m *memoryDB) clean() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items = make(map[string]string)
}

func (m *memoryDB) count() int {
	return len(m.items)
}

func (m *memoryDB) save() {
	f, err := os.Create("./store/db.json")
	if err != nil {
		fmt.Println("не удалось создать файл", err.Error())
	} else if err := json.NewEncoder(f).Encode(m.items); err != nil {
		fmt.Println("не удалось закодировать", err.Error())
	} else {
		fmt.Println("Успешно сохранено ", len(m.items), "записей в файл")
	}
}

func (m *memoryDB) backUp() {
	t := time.Now() //It will return time.Time object with current timestamp

	tUnixMilli := strconv.FormatInt(int64(time.Nanosecond)*t.UnixNano()/int64(time.Millisecond), 10)

	f, err := os.Create("./store/db" + tUnixMilli + ".json")
	if err != nil {
		fmt.Println("не удалось создать файл", err.Error())
	} else if err := json.NewEncoder(f).Encode(m.items); err != nil {
		fmt.Println("не удалось закодировать", err.Error())
	} else {
		fmt.Println("успешное сохранение db в файл")
	}
}
