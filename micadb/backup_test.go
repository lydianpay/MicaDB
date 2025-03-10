package micadb

import (
	"github.com/Tether-Payments/micadb/tests"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestStartBackup(t *testing.T) {
	mica, err := NewMicaDB(
		Options{
			Filename: "./tests/databases/unittest.bin",
			IsTest:   false,
			CustomTypes: []any{
				tests.TestingStruct2{},
				tests.TestingStruct1{},
			},
			BackupFrequency: 1,
		})

	if err != nil {
		t.Errorf("error attempting to load database for creating : %v", err)
	}

	maxInserts := rand.Intn(91) + 10
	items := map[string]any{}

	// Create random items to be stored
	for range maxInserts {
		items[tests.RandomString()] = tests.RandomItem()
	}

	// Store items to in-memory db
	for key, val := range items {
		mica.Set(key, val)
	}

	// Sleep past the inverval
	time.Sleep(2 * time.Second)

	// Check to ensure the file was created
	_, err = os.Stat("./tests/databases/unittest.bin")
	if err != nil {
		t.Errorf("Backup failed to run as expected")
	}
}

func TestBackup(t *testing.T) {

}
