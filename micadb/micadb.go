package micadb

import (
	"encoding/gob"
	"fmt"
	"reflect"
	"slices"
	"sync"
)

var (
	supportedTypes = []string{"int", "[]int", "bool", "[]bool", "float32", "[]float32", "float64", "[]float64", "string", "[]string"}
)

type MicaDB struct {
	mu       sync.Mutex
	KVPs     map[string]any
	Types    []string
	TypesMap map[string]map[string]string
	Options  Options
}

func New(options Options) *MicaDB {
	m := &MicaDB{
		Options: options,
	}

	if m.Options.Filename == "" {
		m.Options.Filename = defaultDBName
	}

	m.KVPs = map[string]any{}
	m.Types = []string{}
	m.TypesMap = map[string]map[string]string{}

	return m
}

func (m *MicaDB) Start() (db *MicaDB, err error) {
	newBackupFrequency := m.Options.BackupFrequency
	err = m.loadLocalDB()
	if err != nil {
		return nil, err
	}

	// Allow for the updating of backup frequency
	if newBackupFrequency != m.Options.BackupFrequency {
		m.Options.BackupFrequency = newBackupFrequency
	}

	// Start the automatic backup if requested
	if m.Options.BackupFrequency > 0 {
		go m.startBackup()
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

// WithCustomTypes is a wrapper for RegisterCustomType to allow for multiple types to be passed
func (m *MicaDB) WithCustomTypes(values ...any) *MicaDB {
	for _, value := range values {
		m.RegisterCustomType(value)
	}
	return m
}

// RegisterCustomType is used to tell MicaDB about different data structures, outside of primitives, that may be
// contained within the database
func (m *MicaDB) RegisterCustomType(value any) *MicaDB {
	gob.Register(value)
	thisType := reflect.TypeOf(value).String()
	if !slices.Contains(supportedTypes, thisType) && !slices.Contains(m.Types, thisType) {
		describeFields(value)
		typeName := reflect.TypeOf(value).String()
		m.Types = append(m.Types, typeName)
		m.TypesMap[typeName] = describeFields(value)
	}

	return m
}

// Set is used to store a key-value pair
func (m *MicaDB) Set(key string, value any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.KVPs[key] = value
}

// Delete is used to remove a key-value pair
func (m *MicaDB) Delete(key string) {
	m.mu.Lock()

	delete(m.KVPs, key)
}

// Get is used to retrieve a value by key
func (m *MicaDB) Get(key string) any {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.KVPs[key]
}

// GetAll is used to retrieve all key-value pairs within MicaDB
func (m *MicaDB) GetAll() map[string]any {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.KVPs
}

// HasKey is used to tell if a given key exists within MicaDB
func (m *MicaDB) HasKey(key string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, x := m.KVPs[key]
	return x
}
