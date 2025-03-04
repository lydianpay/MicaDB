package cmd

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Tether-Payments/micadb/db"
	"github.com/spf13/cobra"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var concurrencyCmd = &cobra.Command{
	Use:   "concurrency",
	Short: "Concurrency test",

	Run: func(cmd *cobra.Command, args []string) {
		db, err := db.NewMicaDB("concurrency.bin")
		if err != nil {
			panic(err)
		}
		wg := sync.WaitGroup{}
		items := buildItems()
		doWrite := func(name string, item any) {
			db.Create(name, item)
			wg.Done()
		}
		wg.Add(quantity)
		pretty := message.NewPrinter(language.English)
		log.Println(pretty.Sprintf("Performing %d concurrent writes", quantity))
		writeNow := time.Now()
		for i := range quantity {
			go doWrite(fmt.Sprintf("Thread #%d", i), items[i])
		}
		wg.Wait()
		log.Println(pretty.Sprintf("Writing %d items took %v", quantity, time.Since(writeNow)))
		allItems := db.GetAll()
		for i := range quantity {
			_, OK := allItems[fmt.Sprintf("Thread #%d", i)]
			if !OK {
				panic("item was missing")
			}
		}
		log.Println(pretty.Sprintf("Recovered %d items", len(allItems)))
	},
}

func buildItems() []any {
	items := []any{}
	for range quantity {
		items = append(items, randomItem())
	}
	return items
}
func init() {
	rootCmd.AddCommand(concurrencyCmd)
	concurrencyCmd.Flags().IntVarP(&quantity, "quantity", "q", quantity, "The quantity of items to insert")
}
