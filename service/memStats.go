package service

import (
	"encoding/json"
	"fmt"
	"runtime"
)

//Monitor структура данных о состоянии памяти
type Monitor struct {
	Alloc,
	TotalAlloc,
	Sys,
	Mallocs,
	Frees,
	LiveObjects,
	PauseTotalNs uint64

	NumGC        uint32
	NumGoroutine int
}

// NewMonitor создаёт экземпляр структуры Monitor с данными на настоящй момент
func NewMonitor() Monitor {
	var m Monitor
	var rtm runtime.MemStats
	//	var interval = time.Duration(duration) * time.Second

	// Read full mem stats
	runtime.ReadMemStats(&rtm)

	// Number of goroutines
	m.NumGoroutine = runtime.NumGoroutine()

	// Misc memory stats
	m.Alloc = rtm.Alloc / 1024
	m.TotalAlloc = rtm.TotalAlloc / 1024
	m.Sys = rtm.Sys / 1024
	m.Mallocs = rtm.Mallocs
	m.Frees = rtm.Frees

	// Live objects = Mallocs - Frees
	m.LiveObjects = m.Mallocs - m.Frees

	// GC Stats
	m.PauseTotalNs = rtm.PauseTotalNs
	m.NumGC = rtm.NumGC

	return m
}

//Get формирует JSON строку для выдачи данных о состоянии памяти
func (mem *Monitor) Get() string {
	// Just encode to json
	b, err := json.Marshal(mem)
	if err != nil {
		fmt.Println("не удалось записать данные о состоянии памяти", err.Error())
	}
	return string(b)
}
