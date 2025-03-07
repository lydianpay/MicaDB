package simulator

import (
	"fmt"
	"github.com/Tether-Payments/micadb/db"
	"log"
	"sync"
	"time"
)

func Concurrency(itemCount int) {
	mica, err := db.NewMicaDB(db.Options{
		Filename: "./tests/databases/concurrency.bin",
		IsTest:   false,
		CustomTypes: []any{
			TestingStruct2{},
			TestingStruct1{},
		},
		BackupFrequency: -1,
	})
	if err != nil {
		panic(err)
	}
	wg := sync.WaitGroup{}
	items := buildItems(itemCount)
	doWrite := func(name string, item any) {
		mica.Set(name, item)
		wg.Done()
	}
	wg.Add(itemCount)
	log.Printf("Performing %d concurrent writes", itemCount)
	writeNow := time.Now()
	for i := range itemCount {
		go doWrite(fmt.Sprintf("Thread #%d", i), items[i])
	}
	wg.Wait()
	log.Printf("Writing %d items took %v", itemCount, time.Since(writeNow))
	allItems := mica.GetAll()
	for i := range itemCount {
		_, OK := allItems[fmt.Sprintf("Thread #%d", i)]
		if !OK {
			log.Fatal("item was missing")
		}
	}
	log.Printf("Recovered %d items", len(allItems))
}

func buildItems(quantity int) []any {
	items := []any{}
	for range quantity {
		items = append(items, RandomItem())
	}
	return items
}
