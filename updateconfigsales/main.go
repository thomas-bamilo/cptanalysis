package main

import (
	"log"
	"sync"
	"time"

	"github.com/globalsign/mgo"
	"github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/bobinteract"
	"github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/mongointeract"
	"github.com/thomas-bamilo/sql/connectdb"
)

func main() {

	start := time.Now()
	log.Println(`Start time: ` + start.Format(`1 January 2006, 15:04:05`))

	dbBob := connectdb.ConnectToBob()
	defer dbBob.Close()
	bamiloCatalogConfigSalesTable := bobinteract.GetBamiloCatalogConfigSalesTable(dbBob)
	bamiloCatalogConfigSalesHistTable := bobinteract.GetBamiloCatalogConfigSalesHistTable(dbBob)

	var url = `mongodb://localhost:27017/competition_analysis`
	mongoSession, err := mgo.Dial(url)
	checkError(err)
	defer mongoSession.Close()

	var wg sync.WaitGroup
	wg.Add(2)
	//wg.Done() inside each function

	// mongo
	go mongointeract.UpsertConfigSales(mongoSession, bamiloCatalogConfigSalesTable, start, &wg)
	go mongointeract.UpsertConfigSalesHist(mongoSession, bamiloCatalogConfigSalesHistTable, start, &wg)

	wg.Wait()

}

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
