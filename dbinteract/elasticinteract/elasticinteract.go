package elasticinteract

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/thomas-bamilo/commercial/competitionanalysis/bamilocatalogconfig"
	"gopkg.in/olivere/elastic.v5"
)

func UpsertConfigInfo(elasticClient *elastic.Client, ctx context.Context, bamiloCatalogConfigTable []bamilocatalogconfig.BamiloCatalogConfig, start time.Time, wg *sync.WaitGroup) {

	defer wg.Done()
	// make it faster: erase all data then re-insert everything, do not provide custom ID, just add IDBmlCatalogConfig as a field
	// and / or use bulk insert (your nemesis)

	_, err := elasticClient.DeleteIndex(`bml_catalog_config`).Do(ctx)
	checkError(err)

	_, err = elasticClient.CreateIndex(`bml_catalog_config`).Do(ctx)
	checkError(err)

	// Setup a bulk processor
	bulkProcessor, err := elasticClient.BulkProcessor().
		Name("bulkProcessor").
		Workers(4).       // number of workers
		BulkActions(950). // commit if # requests >= 950
		Do(ctx)
	if err != nil {
		log.Println(err)
	}

	for _, bamiloCatalogConfig := range bamiloCatalogConfigTable {
		// only keep appropriate information for elastic
		bamiloCatalogConfigElastic := bamilocatalogconfig.BamiloCatalogConfigElastic{
			IDBmlCatalogConfig: bamiloCatalogConfig.IDBmlCatalogConfig,
			SKUName:            bamiloCatalogConfig.SKUName,
			Description:        bamiloCatalogConfig.Description,
			ShortDescription:   bamiloCatalogConfig.ShortDescription,
			PackageContent:     bamiloCatalogConfig.PackageContent,
			ProductWarranty:    bamiloCatalogConfig.ProductWarranty,

			BiCategoryOneName:   bamiloCatalogConfig.BiCategoryOneName,
			BiCategoryTwoName:   bamiloCatalogConfig.BiCategoryTwoName,
			BiCategoryThreeName: bamiloCatalogConfig.BiCategoryThreeName,

			BrandName:   bamiloCatalogConfig.BrandName,
			BrandNameEn: bamiloCatalogConfig.BrandNameEn,

			SupplierName:   bamiloCatalogConfig.SupplierName,
			SupplierNameEn: bamiloCatalogConfig.SupplierNameEn,

			AvgPrice:        bamiloCatalogConfig.AvgPrice,
			AvgSpecialPrice: bamiloCatalogConfig.AvgSpecialPrice,
		}

		// convert information to JSON
		//bamiloCatalogConfigElasticByte, err := json.Marshal(bamiloCatalogConfigElastic)
		//checkError(err)
		//bamiloCatalogConfigElasticJSON := string(bamiloCatalogConfigElasticByte)

		bulkIndexRequest := elastic.NewBulkIndexRequest().
			Index("bml_catalog_config").
			Type("bml_catalog_config").
			Doc(bamiloCatalogConfigElastic)

		bulkProcessor.Add(bulkIndexRequest)

		/*// index JSON information to JSON
		_, err = elasticClient.Index().
			Index(`bml_catalog_config`).
			Type(`bml_catalog_config`).
			//Id(strconv.Itoa(bamiloCatalogConfig.IDBmlCatalogConfig)).
			BodyJson(bamiloCatalogConfigElasticJSON).
			Do(ctx)

		checkError(err)*/
	}

	end := time.Now()
	log.Println(`End time Elastic: ` + end.Format(`1 January 2006, 15:04:05`))
	duration := time.Since(start)
	log.Print(`Time elapsed Elastic: `, duration.Minutes(), ` minutes`)

	err = bulkProcessor.Close()
	checkError(err)

}

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
