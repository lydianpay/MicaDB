package filesystem

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const Separator = `⌘`

type fsdb struct {
	kvps     map[string]any
	filename string
}

func NewFilesystemDB(filename string) (*fsdb, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error loading file '%s' : %v", filename, err)
	}
	defer file.Close()
	kvps := map[string]any{}
	fileScanner := bufio.NewScanner(file)
	for fileScanner.Scan() {
		record := fileScanner.Text()
		split := strings.Split(record, Separator)
		if len(split) == 2 {
			kvps[split[0]] = split[1]
		}
	}

	return &fsdb{kvps: kvps, filename: filename}, nil
}

func (f *fsdb) Create(key string, value any) {
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr:
		fmt.Println("Got an int-y thing")
		valStr := fmt.Sprintf("%v", value)
		valInt, err := strconv.Atoi(valStr)
		if err != nil {
			return
		}
		fmt.Printf("Maybe seperate like this [%v]%v\n", reflect.TypeOf(value), value)
		fmt.Printf("Or maybe %s%s[%vb64]%v\n", key, Separator, reflect.TypeOf(value), base64.StdEncoding.EncodeToString(big.NewInt(int64(valInt)).Bytes()))
	// case uint, uint8, uint16, uint32, uint64, uintptr:
	// 	fmt.Println("Got an unsigned int-y thing")
	case float32, float64:
		fmt.Println("Got a float-y thing")
	case complex64, complex128:
		fmt.Println("Got a complex thing")
	case bool:
		fmt.Println("Got a bool")
	case string:
		fmt.Println("Got a string")
	case interface{}:
		jb, err := json.Marshal(value)
		if err != nil {
			fmt.Println("got something i couldn't json-ify")
		}
		fmt.Printf("got something else and it looks like this : %s\n", jb)
	default:
		fmt.Println("I don't know what it is")
	}
	// fmt.Println(reflect.TypeOf(value))
	f.kvps[key] = value
}

func (f *fsdb) HasKey(key string) bool {
	_, x := f.kvps[key]
	return x
}

func (f *fsdb) Retrieve(key string) any {
	return f.kvps[key]
}
func (f *fsdb) GetAll() map[string]any {
	return f.kvps
}
func (f *fsdb) Update(key string, value any) {
	f.kvps[key] = value
}

func (f *fsdb) Delete(key string) {
	delete(f.kvps, key)
}

func (f *fsdb) Save() error {
	os.Rename(f.filename, fmt.Sprintf("%s.%s.%s", f.filename, time.Now().Format(time.RFC3339), "backup"))
	file, err := os.OpenFile(f.filename, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error opening file '%s' for writing:%v", f.filename, err)
	}
	defer file.Close()
	for key, val := range f.kvps {
		_, err = file.Write([]byte(fmt.Sprintf("%s%s%v\n", key, Separator, val)))
		if err != nil {
			return fmt.Errorf("error writing the record '%s' to the file '%s' (value redacted) : %v", key, f.filename, err)
		}
	}

	return nil
}
