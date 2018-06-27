package main

import (
	"log"
	"time"

	"github.com/thomas-bamilo/sql/connectdb"

	"github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/baainteract"
	"github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/bobinteract"
)

func main() {

	start := time.Now()
	log.Println(`Start time: ` + start.Format(`1 January 2006, 15:04:05`))

	dbBob := connectdb.ConnectToBob()
	defer dbBob.Close()
	bamiloCatalogConfigSalesTable := bobinteract.GetBamiloCatalogConfigSalesTable(dbBob)

	dbBaa := connectdb.ConnectToBaa()
	defer dbBaa.Close()
	baainteract.UpdateBamiloCatalogConfigSales(dbBaa, bamiloCatalogConfigSalesTable)

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
