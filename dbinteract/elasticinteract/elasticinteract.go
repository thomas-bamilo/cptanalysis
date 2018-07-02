package elasticinteract

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/thomas-bamilo/commercial/competitionanalysis/bamilocatalogconfig"
	"gopkg.in/olivere/elastic.v5"
)

func UpsertConfigInfo(elasticClient *elastic.Client, ctx context.Context, bamiloCatalogConfigTable []bamilocatalogconfig.BamiloCatalogConfig, start time.Time) {

	for _, bamiloCatalogConfig := range bamiloCatalogConfigTable {
		// only keep appropriate information for elastic
		bamiloCatalogConfigElastic := bamilocatalogconfig.BamiloCatalogConfig{
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
		bamiloCatalogConfigElasticByte, err := json.Marshal(bamiloCatalogConfigElastic)
		checkError(err)
		bamiloCatalogConfigElasticJSON := string(bamiloCatalogConfigElasticByte)

		// index JSON information to JSON
		_, err = elasticClient.Index().
			Index(`bamilo_catalog_config`).
			Type(`bamilo_catalog_config`).
			Id(bamiloCatalogConfig.IDBmlCatalogConfig).
			BodyJson(bamiloCatalogConfigElasticJSON).
			Do(ctx)

		checkError(err)
	}

	end := time.Now()
	log.Println(`End time Elastic: ` + end.Format(`1 January 2006, 15:04:05`))
	duration := time.Since(start)
	log.Print(`Time elapsed Elastic: `, duration.Minutes(), ` minutes`)

}

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
