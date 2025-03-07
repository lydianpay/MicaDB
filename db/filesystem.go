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
		log.Fatal(err)
	}

	return nil
}
