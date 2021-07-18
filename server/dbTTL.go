package server

import (
	"fmt"
	"strconv"
	"sync"
)

type timeToLive struct {
	mu     sync.RWMutex
	keyTTL map[string]int64
}

//newTTLdb инициализирует базу TTL меток
func newTTLdb() timeToLive {
	return timeToLive{keyTTL: map[string]int64{}}
}

//setTTL метод установки времени жизни
func (t *timeToLive) setTTL(key string, value string) {
	t.mu.Lock()

	time, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		fmt.Println("не удалось установить TTL", err.Error())
	}
	t.keyTTL[key] = time
	t.mu.Unlock()
}

//ttlKill убивает по истечению времени