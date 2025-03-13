package micadb

import (
	"errors"
	"github.com/Tether-Payments/micadb/tests"
	"math/rand"
	"os"
	"testing"
)

func TestCreatePersistentStorage(t *testing.T) {
	db := createTestDB(UnitTestDBFilename)

	_, err := db.createPersistentStorage()
	if err != nil {
		t.Error(err)
	}

	// File should exist
	_, err = os.Stat(UnitTestDBFilename)
	if err != nil {
		t.Error(err)
	}

	cleanupTestFiles()

	db = createTestDB(UnitTestDBBadFilename)
	_, err = db.createPersistentStorage()
	if err == nil {
		t.Error("This test should have failed")
	}

	cleanupTestFiles()
}

func TestLoadLocalDB(t *testing.T) {
	_, err := createReadableTestDB()
	if err != nil {
		t.Error(err)
	}

	dbReader := createTestDB(UnitTestDBFilename)
	err = dbReader.loadLocalDB()
	if err != nil {
		t.Error(err)
	}

	if dbReader.Get(UnitTestKey) != UnitTestVal {
		t.Error("Expected bar, got ", dbReader.Get("foo"))
	}

	cleanupTestFiles()
}

func TestFailWithTypes(t *testing.T) {

	dbNoRegistration, err := createReadableTestDB()
	if err != nil {
		t.Error(err)
	}

	// Add a value
	dbNoRegistration.Set("forcedCustom", tests.TestingStruct3{
		AnotherString: "Test",
		AnotherInt:    7,
	})

	err = dbNoRegistration.Backup()

	db := createTestDB(UnitTestDBFilename)

	// Start MicaDB without registering custom types
	_, err = db.Start()
	if errors.Is(ErrDecodingFileWithTypes, err) != true {
		t.Errorf("Expected decoding error with file types, received: %v", err)
	}

	// Remove the .types file to force the default error
	_ = os.Remove(UnitTestDBFilename + ".types")

	_, err = db.Start()
	if errors.Is(ErrTypesFileNotFound, err) != true {
		t.Errorf("Expected decoding error without file types, received: %v", err)
	}

	cleanupTestFiles()
}

func createTestDB(filename string) *MicaDB {
	return New(Options{
		Filename:        filename,
		IsTest:          false,
		BackupFrequency: -1,
	})
}

func createReadableTestDB() (db *MicaDB, err error) {
	db, err = New(
		Options{
			Filename:        UnitTestDBFilename,
			IsTest:          false,
			BackupFrequency: 1,
		}).WithCustomTypes(
		tests.TestingStruct2{},
		tests.TestingStruct1{},
	).Start()

	if err != nil {
		return nil, err
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

	db.Set(UnitTestKey, UnitTestVal)

	err = db.Backup()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func cleanupTestFiles() {
	_ = os.Remove(UnitTestDBFilename)
	_ = os.Remove(UnitTestDBFilename + ".types")
}
