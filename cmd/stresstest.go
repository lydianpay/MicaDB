/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/Tether-Payments/micadb/db"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// stresstestCmd represents the stresstest command
var stresstestCmd = &cobra.Command{
	Use:   "stresstest",
	Short: "A brief description of your command",

	Run: func(cmd *cobra.Command, args []string) {
		_, err := os.Stat("./stresstest.lastphase")
		if err == nil {
			lpb, err := os.ReadFile("./stresstest.lastphase")
			if err != nil {
				panic(err)
			}
			if string(lpb) == "read" {
				log.Println("Performing Write test")
				writeTest()
				fmt.Println("\n\nWrite test performed. Perform stresstest again for read test")
			} else {
				log.Println("Performing Read test")
				readTest()
			}
		} else {
			log.Println("Performing Write test")
			writeTest()
			fmt.Println("\n\nWrite test performed. Perform stresstest again for read test")
		}

	},
}

func readTest() {
	log.Println("Starting Read")
	starTime := time.Now()
	db, err := db.NewMicaDB("./stresstest.bin", TestingStruct2{}, TestingStruct1{})
	if err != nil {
		log.Fatalf("error attempting to load database for creating : %v", err)
	}
	loadDuration := time.Since(starTime)
	getAllTime := time.Now()
	all := db.GetAll()
	for key, val := range all {
		_ = key
		_ = val
	}
	readDuration := time.Since(getAllTime)
	p := message.NewPrinter(language.English)
	log.Println(p.Sprintf("Items loaded : %v ", len(all)))
	log.Printf("Load duration : %v", loadDuration)
	log.Printf("Read duration : %v", readDuration)
	os.WriteFile("./stresstest.lastphase", []byte("read"), 0655)
}

func writeTest() {
	_, err := os.Stat("./stresstest.bin")
	if err == nil {
		fmt.Println("deleting existing ./stresstest.bin")
		rErr := os.Remove("./stresstest.bin")
		if rErr != nil {
			log.Fatalf("couldn't remove existing stresstest.bin : %v", err)
		}
	}
	db, err := db.NewMicaDB("./stresstest.bin")
	if err != nil {
		log.Fatalf("error attempting to load database for creating : %v", err)
	}
	pretty := message.NewPrinter(language.English)
	maxInserts := quantity
	items := map[string]any{}

	// Create random items to be stored
	log.Println(pretty.Sprintf("Creating %v random items", maxInserts))
	randomItemsNow := time.Now()
	for range maxInserts {
		items[rs()] = randomItem()
	}
	log.Printf("Random item creation took %v", time.Since(randomItemsNow))

	// Store items to in-memory db
	log.Println("Starting Memory Write")
	memWriteNow := time.Now()
	for key, val := range items {
		db.Create(key, val)
	}
	log.Printf("Writing items to in-memory database took %v", time.Since(memWriteNow))

	// Store db to file
	log.Println("Saving to file")
	fileSaveNow := time.Now()
	err = db.Save()
	if err != nil {
		panic(err)
	}
	log.Printf("Writing database to file took %v", time.Since(fileSaveNow))

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	log.Println(pretty.Sprintf("Record Count : %v ", maxInserts))
	log.Println(pretty.Sprintf("Alloc       : %v bytes", memStats.Alloc))
	log.Println(pretty.Sprintf("Total Alloc : %v bytes", memStats.TotalAlloc))
	log.Println(pretty.Sprintf("Sys         : %v bytes", memStats.Sys))
	log.Println(pretty.Sprintf("GC Count    : %v", memStats.NumGC))
	fi, err := os.Stat("./stresstest.bin")
	if err != nil {
		panic(err)
	} else {
		log.Println(pretty.Sprintf("File size : %v bytes (approx. %v MB)", fi.Size(), fi.Size()/1024/1024))
	}
	log.Println("Write test Done")
	os.WriteFile("./stresstest.lastphase", []byte("write"), 0655)
}

func init() {
	rootCmd.AddCommand(stresstestCmd)
	stresstestCmd.Flags().IntVarP(&quantity, "qty", "q", quantity, "The number of random inserts to perform")
}
