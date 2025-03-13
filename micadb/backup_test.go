package micadb

import (
	"github.com/Tether-Payments/micadb/tests"
	"math/rand"
	"os"
	"testing"
	"time"
)

const (
	UnitTestDBFilename    = "./../tests/databases/unittest.bin"
	UnitTestDBBadFilename = "./nonexistent/directory/unittest.bin"
	UnitTestKey           = "foo"
	UnitTestVal           = "bar"
)

func TestStartBackup(t *testing.T) {
	db, err := New(
		Options{
			Filename:        UnitTestDBFilename,
			IsTest:          false,
			BackupFrequency: 1,
		}).WithCustomTypes(
		tests.TestingStruct2{},
		tests.TestingStruct1{},
	).Start()

	if err != nil {
		t.Errorf("error attempting to load database : %v", err)
	}

	maxInserts := rand.Intn(91) + 10
	items := map[string]any{}

	// Create random items to be stored
	for range maxInserts {
		items[tests.RandomString()] = tests.RandomItem()
	}

	// Store items to in-memory db
	for key, val := range items {
		db.Set(key, val)
	}

	// Sleep past the inverval
	time.Sleep(2 * time.Second)

	// Check to ensure the file was created
	_, err = os.Stat(UnitTestDBFilename)
	if err != nil {
		t.Errorf("Backup failed to run as expected")
	}

	cleanupTestFiles()
}

func TestBackup(t *testing.T) {
	_, err := New(
		Options{
			Filename:        UnitTestDBBadFilename,
			IsTest:          false,
			BackupFrequency: 0,
		}).WithCustomTypes(
		tests.TestingStruct2{},
		tests.TestingStruct1{},
	).Start()

	if err != nil {
		t.Errorf("error attempting to load database : %v", err)
	}

	cleanupTestFiles()
}
