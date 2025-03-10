package db

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func (m *MicaDB) startBackup() {

	ticker := time.NewTicker(time.Duration(m.Options.BackupFrequency) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		m.Backup()
	}
}

func (m *MicaDB) Backup() {
	m.mu.Lock()
	defer m.mu.Unlock()

	file, err := m.CreatePersistentStorage()
	if err != nil {
		log.Printf("error creating database file '%s' : %v", m.Options.Filename, err)
	}
	defer file.Close()

	bufferedWriter := bufio.NewWriter(file)

	gob.Register(m)

	encoder := gob.NewEncoder(bufferedWriter)
	err = encoder.Encode(m)
	if err != nil {
		log.Printf("error encoding database : %v", err)
	}

	err = bufferedWriter.Flush()
	if err != nil {
		log.Printf("error flushing buffered writer : %v", err)
	}

	err = m.persistDataTypes()
	if err != nil {
		log.Printf("error persisting database types : %v", err)
	}
}

func (m *MicaDB) persistDataTypes() error {
	sb2 := strings.Builder{}
	for dataType, fieldDescriptions := range m.TypesMap {
		sb2.WriteString(fmt.Sprintf("[%s]\n", dataType))
		for fieldName, fieldType := range fieldDescriptions {
			sb2.WriteString(fmt.Sprintf("    [%s] : %s\n", fieldName, fieldType))
		}
	}
	err := os.WriteFile(m.Options.Filename+".types", []byte(sb2.String()), 0655)
	if err != nil {
		return fmt.Errorf("couldn't write .types file : %v", err)
	}

	return nil
}
