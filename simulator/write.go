package simulator

import (
	"github.com/Tether-Payments/micadb/db"
	"log"
	"os"
	"runtime"
	"time"
)

func Write(runCount int) {
	mica, err := db.NewMicaDB(
		db.Options{
			Filename: "./tests/databases/stresstest.bin",
			IsTest:   false,
			CustomTypes: []any{
				TestingStruct2{},
				TestingStruct1{},
			},
			BackupFrequency: 0,
		})

	if err != nil {
		log.Fatalf("error attempting to load database for creating : %v", err)
	}

	maxInserts := runCount
	items := map[string]any{}

	// Create random items to be stored
	log.Printf("Creating %v random items", maxInserts)
	randomItemsNow := time.Now()
	for range maxInserts {
		items[RandomString()] = RandomItem()
	}
	log.Printf("Random item creation took %v", time.Since(randomItemsNow))

	// Store items to in-memory db
	log.Println("Starting Memory Write")
	memWriteNow := time.Now()
	for key, val := range items {
		mica.Set(key, val)
	}
	log.Printf("Writing %d items to in-memory database took %v", len(items), time.Since(memWriteNow))

	// Store db to file
	log.Println("Saving to file")
	fileSaveNow := time.Now()
	mica.Backup()

	log.Printf("Writing database to file took %v", time.Since(fileSaveNow))

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	log.Printf("Record Count : %v ", maxInserts)
	log.Printf("Alloc       : %v bytes", memStats.Alloc)
	log.Printf("Total Alloc : %v bytes", memStats.TotalAlloc)
	log.Printf("Sys         : %v bytes", memStats.Sys)
	log.Printf("GC Count    : %v", memStats.NumGC)
	fi, err := os.Stat("./tests/databases/stresstest.bin")
	if err != nil {
		panic(err)
	} else {
		log.Printf("File size : %v bytes (approx. %v MB)", fi.Size(), fi.Size()/1024/1024)
	}
	log.Println("Write test Done")
	os.WriteFile("./stresstest.lastphase", []byte("write"), 0655)
}
