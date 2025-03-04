package db

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"reflect"
	"slices"
	"strings"
	"sync"
)

var (
	ignoretypes = []string{"int", "[]int", "bool", "[]bool", "float32", "[]float32", "float64", "[]float64", "string", "[]string"}
)

type MicaDB struct {
	mu       sync.Mutex
	KVPs     map[string]any
	Filename string
	Types    []string
	TypesMap map[string]map[string]string
}

func NewMicaDB(filename string, sampleItems ...any) (*MicaDB, error) {
	m := &MicaDB{}
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		log.Printf("DB '%s' does not exist. Creating", filename)
		m.Filename = filename
		m.KVPs = map[string]any{}
		m.Types = []string{}
		m.TypesMap = map[string]map[string]string{}
	} else {
		file, err := os.Open(filename)
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
		defer file.Close()
		gob.Register(m)
		for _, sample := range sampleItems {
			gob.Register(sample)
		}
		decoder := gob.NewDecoder(file)
		err = decoder.Decode(m)
		if err != nil {
			statFile, statErr := os.ReadFile(filename + ".types")
			if statErr != nil {
				return nil, fmt.Errorf("error decoding db file contents and also an error loading the associated .types file : %v : %v", err, statErr)
			}
			if statFile != nil {
				return nil, fmt.Errorf("error decoding db file contents. Make sure you initialize the db with each of the non-primitives in this list : \n%s\nOriginal error : %v", statFile, err)
			} else {
				return nil, fmt.Errorf("error decoding db file contents : %v", err)
			}
		}

	}
	return m, nil

}

func describeFields(value any) map[string]string {
	descriptions := map[string]string{}
	fields := reflect.VisibleFields(reflect.TypeOf(value))
	for _, field := range fields {
		descriptions[fmt.Sprintf("%v", field.Name)] = fmt.Sprintf("%v", field.Type)
	}
	return descriptions
}

func (m *MicaDB) registerTheseTypes(values ...any) *MicaDB {
	for _, value := range values {
		gob.Register(value)
		thisType := reflect.TypeOf(value).String()
		if !slices.Contains(ignoretypes, thisType) && !slices.Contains(m.Types, thisType) {
			describeFields(value)
			typeName := reflect.TypeOf(value).String()
			m.Types = append(m.Types, typeName)
			m.TypesMap[typeName] = describeFields(value)
		}
	}
	return m
}

func (m *MicaDB) Create(key string, value any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.registerTheseTypes(value)
	m.KVPs[key] = value
}

func (m *MicaDB) HasKey(key string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, x := m.KVPs[key]
	return x
}

func (m *MicaDB) Retrieve(key string) any {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.KVPs[key]
}

func (m *MicaDB) GetAll() map[string]any {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.KVPs
}

func (m *MicaDB) Update(key string, value any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.KVPs[key] = value
}

func (m *MicaDB) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.KVPs, key)
}

func (m *MicaDB) Save() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	file, err := os.Create(m.Filename)
	if err != nil {
		return fmt.Errorf("error creating database file '%s' : %v", m.Filename, err)
	}
	defer file.Close()
	bufferedWriter := bufio.NewWriter(file)
	gob.Register(m)
	encoder := gob.NewEncoder(bufferedWriter)
	err = encoder.Encode(m)
	if err != nil {
		return fmt.Errorf("error encoding database : %v", err)
	}
	err = bufferedWriter.Flush()
	if err != nil {
		return fmt.Errorf("error flushing buffered writer : %v", err)
	}
	sb2 := strings.Builder{}
	for dataType, fieldDescriptions := range m.TypesMap {
		sb2.WriteString(fmt.Sprintf("[%s]\n", dataType))
		for fieldName, fieldType := range fieldDescriptions {
			sb2.WriteString(fmt.Sprintf("    [%s] : %s\n", fieldName, fieldType))
		}
	}
	err = os.WriteFile(m.Filename+".types", []byte(sb2.String()), 0655)
	if err != nil {
		return fmt.Errorf("couldn't write .types file : %v", err)
	}
	return nil
}
