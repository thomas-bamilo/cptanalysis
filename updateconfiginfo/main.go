package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/globalsign/mgo"
	"github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/bobinteract"
	"github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/elasticinteract"
	"github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/mongointeract"
	"github.com/thomas-bamilo/sql/connectdb"
	elastic "gopkg.in/olivere/elastic.v5"
)

func main() {

	start := time.Now()
	log.Println(`Start time: ` + start.Format(`1 January 2006, 15:04:05`))

	dbBob := connectdb.ConnectToBob()
	defer dbBob.Close()
	bamiloCatalogConfigTable := bobinteract.GetBamiloCatalogConfigTable(dbBob)

	// Connection URL
	var url = `mongodb://localhost:27017/competition_analysis`
	mongoSession, err := mgo.Dial(url)
	checkError(err)
	defer mongoSession.Close()

	elasticClient, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"))
	checkError(err)
	ctx := context.Background()

	// wg allows the program to wait for both goroutines to execute before ending
	// otherwise, the program will end before any goroutine can finish!
	var wg sync.WaitGroup
	wg.Add(3)

	go mongointeract.UpsertConfigInfo(mongoSession, bamiloCatalogConfigTable, start)

	go mongointeract.UpsertConfigHistory(mongoSession, bamiloCatalogConfigTable, start)

	go elasticinteract.UpsertConfigInfo(elasticClient, ctx, bamiloCatalogConfigTable, start)

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
