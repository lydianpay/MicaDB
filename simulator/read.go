package simulator

import (
	"github.com/Tether-Payments/micadb/micadb"
	"github.com/Tether-Payments/micadb/tests"
	"log"
	"os"
	"time"
)

func Read() {
	log.Println("Starting Read")
	starTime := time.Now()

	db, err := micadb.New(micadb.Options{
		Filename:        "./tests/databases/stresstest.bin",
		IsTest:          false,
		BackupFrequency: -1,
	}).WithCustomTypes(
		tests.TestingStruct2{},
		tests.TestingStruct1{},
	).Start()

	if err != nil {
		log.Fatalf("error attempting to load database for creating : %v", err)
	}

	loadDuration := time.Since(starTime)
	getAllTime := time.Now()

	all := db.GetAll()
	for key, val := range all {
		_ = key
		_ = val
	}

	readDuration := time.Since(getAllTime)

	log.Printf("Items loaded : %v ", len(all))
	log.Printf("Load duration : %v", loadDuration)
	log.Printf("Read duration : %v", readDuration)
	os.WriteFile("./stresstest.lastphase", []byte("read"), 0655)
}
