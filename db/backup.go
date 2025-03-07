package db

import (
	"bufio"
	"encoding/gob"
	"log"
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
}
