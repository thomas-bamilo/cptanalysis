package main

import (
	"context"
	"log"
	"time"

	"github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/bobinteract"
	"github.com/thomas-bamilo/sql/connectdb"
	elastic "gopkg.in/olivere/elastic.v5"
)

func main() {

	start := time.Now()
	log.Println(`Start time: ` + start.Format(`1 January 2006, 15:04:05`))

	elasticClient, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"))
	checkError(err)
	ctx := context.Background()

	dbBob := connectdb.ConnectToBob()
	defer dbBob.Close()

	dbBaa := connectdb.ConnectToBaa()
	defer dbBaa.Close()

	bobinteract.BamiloCatalogConfigFromBobToElasticAndBaa(dbBob, elasticClient, dbBaa, ctx)

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
