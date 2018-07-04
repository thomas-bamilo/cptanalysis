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
	mongoSession2 := mongoSession.Copy()
	defer mongoSession2.Close()

	elasticClient, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"))
	checkError(err)
	ctx := context.Background()

	// wg allows the program to wait for goroutines to execute before ending
	// otherwise, the program will end before any goroutine can finish!
	var wg sync.WaitGroup
	wg.Add(3)

	// mongo
	// snapshot run twice per day
	go mongointeract.UpsertConfigInfo(mongoSession, bamiloCatalogConfigTable, start, &wg)
	// history only run once per day
	go mongointeract.UpsertConfigInfoHist(mongoSession2, bamiloCatalogConfigTable, start, &wg)
	// elastic
	go elasticinteract.UpsertConfigInfo(elasticClient, ctx, bamiloCatalogConfigTable, start, &wg)

	wg.Wait()

}

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
