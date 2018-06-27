package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/baainteract"
	"github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/bobinteract"
	"github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/elasticinteract"
	"github.com/thomas-bamilo/sql/connectdb"
	elastic "gopkg.in/olivere/elastic.v5"
)

func main() {

	start := time.Now()
	log.Println(`Start time: ` + start.Format(`1 January 2006, 15:04:05`))

	dbBob := connectdb.ConnectToBob()
	defer dbBob.Close()
	bamiloCatalogConfigTable := bobinteract.GetBamiloCatalogConfigTable(dbBob)

	dbBaa := connectdb.ConnectToBaa()
	defer dbBaa.Close()

	elasticClient, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"))
	checkError(err)
	ctx := context.Background()
	// wg allows the program to wait for both goroutines to execute before ending
	// otherwise, the program will end before any goroutine can finish!
	var wg sync.WaitGroup
	wg.Add(2)

	go baainteract.AddBamiloCatalogConfigTableToBaa(dbBaa, bamiloCatalogConfigTable, start)

	go elasticinteract.AddBamiloCatalogConfigTableToIndex(elasticClient, ctx, bamiloCatalogConfigTable, start)

	wg.Wait()

	end := time.Now()
	log.Println(`End time: ` + end.Format(`1 January 2006, 15:04:05`))
	duration := time.Since(start)
	log.Print(`Time elapsed: `, duration.Minutes(), ` minutes`)

}

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
