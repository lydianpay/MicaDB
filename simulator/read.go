package simulator

import (
	"log"
	"os"
	"time"

	"github.com/tetherpay/micadb/micadb"
	"github.com/tetherpay/micadb/tests"
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
	os.WriteFile("./tests/databases/stresstest.lastphase", []byte("read"), 0655)
}
