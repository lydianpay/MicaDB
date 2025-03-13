package micadb

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var backupTicker *time.Ticker

func (m *MicaDB) startBackup() {

	backupTicker = time.NewTicker(time.Duration(m.Options.BackupFrequency) * time.Second)
	defer backupTicker.Stop()

	for range backupTicker.C {
		err := m.Backup()
		if err != nil {
			log.Println("Backup failed:", err)
		}
	}
}

func (m *MicaDB) Backup() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	file, err := m.createPersistentStorage()
	if err != nil {
		return err
	}

	defer file.Close()

	bufferedWriter := bufio.NewWriter(file)

	gob.Register(m)

	encoder := gob.NewEncoder(bufferedWriter)
	err = encoder.Encode(m)
	if err != nil {
		return err
	}

	err = bufferedWriter.Flush()
	if err != nil {
		return err
	}

	err = m.persistDataTypes()
	if err != nil {
		return err
	}

	return nil
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
