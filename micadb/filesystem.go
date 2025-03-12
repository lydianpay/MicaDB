package micadb

import (
	"encoding/gob"
	"errors"
	"os"
)

var ErrTypesFileNotFound = errors.New(".types file not found and mica db file could not be decoded")
var ErrDecodingFileWithTypes = errors.New("non-primitives found in mica db file, check the .types file")

func (m *MicaDB) createPersistentStorage() (dbFile *os.File, err error) {
	if m.Options.IsTest {
		dbFile, err = os.CreateTemp("", m.Options.Filename)
	} else {
		dbFile, err = os.Create(m.Options.Filename)
	}

	return dbFile, err
}

func (m *MicaDB) loadLocalDB() error {
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
		return err
	}

	defer file.Close()

	// Register the MicaDB structure
	gob.Register(m)

	// Decode the persistent storage db file
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(m)
	if err != nil {
		return m.failWithTypes(err)
	}

	return nil
}

// failWithTypes attempts to read the .types file to sugar the error
func (m *MicaDB) failWithTypes(err error) error {

	_, statErr := os.ReadFile(m.Options.Filename + ".types")
	if statErr != nil {
		return ErrTypesFileNotFound
	}

	return ErrDecodingFileWithTypes
}
