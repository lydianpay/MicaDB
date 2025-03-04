package main

import (
	"log"
	"math/rand"

	tpgstrings "github.com/Tether-Payments/go-shared/strings"
	"github.com/Tether-Payments/micadb/filesystem"
)

type DumbThing struct {
	AString string
	AnInt   int
}

func main() {
	db, err := filesystem.NewFilesystemDB("./somerecords.txt")
	if err != nil {
		panic(err)
	}
	db.Create(tpgstrings.GenerateRandomString(rand.Intn(100), tpgstrings.AlphaNumericWithSymbols_CS), tpgstrings.GenerateRandomString(rand.Intn(100), tpgstrings.AlphaNumericWithSymbols_CS))
	db.Create(tpgstrings.GenerateRandomString(rand.Intn(100), tpgstrings.AlphaNumericWithSymbols_CS), rand.Intn(500000000000))
	db.Create(tpgstrings.GenerateRandomString(rand.Intn(100), tpgstrings.AlphaNumericWithSymbols_CS), rand.Float32()*100)
	db.Create(tpgstrings.GenerateRandomString(rand.Intn(100), tpgstrings.AlphaNumericWithSymbols_CS), rand.Float64()*100)
	db.Create(tpgstrings.GenerateRandomString(rand.Intn(100), tpgstrings.AlphaNumericWithSymbols_CS), false)
	db.Create(tpgstrings.GenerateRandomString(rand.Intn(100), tpgstrings.AlphaNumericWithSymbols_CS), DumbThing{AString: "a", AnInt: 1})
	db.Create(tpgstrings.GenerateRandomString(rand.Intn(100), tpgstrings.AlphaNumericWithSymbols_CS), &DumbThing{AString: "b", AnInt: 2})

	err = db.Save()
	if err != nil {
		log.Printf("bad save : %v", err)
	}
	// records := db.GetAll()
	// for key, val := range records {
	// 	fmt.Printf("Key : %s\n", key)
	// 	fmt.Printf("Val : %v\n", val)
	// 	fmt.Println("---")
	// }
}
