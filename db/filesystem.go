package db

import (
	"encoding/gob"
	"fmt"
	"log"
	"os"
)

func (m *MicaDB) CreatePersistentStorage() (dbFile *os.File, err error) {
	if m.Options.IsTest {
		fmt.Println("Using test mode")
		dbFile, err = os.CreateTemp("", m.Options.Filename)
	} else {
		dbFile, err = os.Create(m.Options.Filename)
	}

	return dbFile, err
}

func (m *MicaDB) LoadLocalDB() error {
	// Check for an existing db file
	_, err := os.Stat(m.Options.Filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		} else {
			return err
		}
	}

	file, err := os.Open(m.Options.Filename)
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}

	defer file.Close()

	// Register the MicaDB structure
	gob.Register(m)

	// Decode the persistent storage db file
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(m)
	if err != nil {
		trueErr := m.failWithTypes(err)
		log.Fatal(trueErr)
	}

	return nil
}

// failWithTypes attempts to read the .types file to sugar the error
func (m *MicaDB) failWithTypes(err error) error {

	statFile, statErr := os.ReadFile(m.Options.Filename + ".types")
	if statErr != nil {
		return fmt.Errorf("error decoding db file contents and also an error loading the associated .types file : %v : %v", err, statErr)
	}
	if statFile != nil {
		return fmt.Errorf("error decoding db file contents. Make sure you initialize the db with each of the non-primitives in this list : \n%s\nOriginal error : %v", statFile, err)
	}

	return fmt.Errorf("error decoding db file contents : %v", err)
}
