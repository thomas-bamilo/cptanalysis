package mongointeract

import (
	"log"
	"time"

	"github.com/globalsign/mgo"
	"github.com/thomas-bamilo/commercial/competitionanalysis/bamilocatalogconfig"
)

func UpsertConfigInfo(mongoSession *mgo.Session, bamiloCatalogConfigTable []bamilocatalogconfig.BamiloCatalogConfig, start time.Time) {

	bamiloCatalogConfigCollection := mongoSession.DB(`competition_analysis`).C(`bml_catalog_config`)

	mybulk := bamiloCatalogConfigCollection.Bulk()
	mybulk.Unordered()
	ops := 0
	now := time.Now()

	for _, bamiloCatalogConfig := range bamiloCatalogConfigTable {

		if bamiloCatalogConfig.VisibleInShop == `1` {
			bamiloCatalogConfig.SetVisibleInShopTrue()
		}

		// only keep appropriate information for bamilocatalogconfig
		bamiloCatalogConfigInfo := bamilocatalogconfig.BamiloCatalogConfig{
			IDMongoDb:          bamiloCatalogConfig.IDBmlCatalogConfig,
			IDBmlCatalogConfig: bamiloCatalogConfig.IDBmlCatalogConfig,
			ConfigSnapshotAt:   now,
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

			VisibleInShopBool: bamiloCatalogConfig.VisibleInShopBool,

			AvgPrice:           bamiloCatalogConfig.AvgPrice,
			AvgSpecialPrice:    bamiloCatalogConfig.AvgSpecialPrice,
			SumOfStockQuantity: bamiloCatalogConfig.SumOfStockQuantity,
			MinOfStockQuantity: bamiloCatalogConfig.MinOfStockQuantity,
		}

		var data interface{}
		data = bamiloCatalogConfigInfo

		if ops > 950 {
			// We need to run our operations before we queue up too many
			_, err := mybulk.Run()
			checkError(err)

			// re-initialize the bulk
			mybulk = nil
			mybulk := bamiloCatalogConfigCollection.Bulk()
			mybulk.Unordered()

			// re-initialize ops
			ops = 0
		}

		ops++
		mybulk.Upsert(data)
	}
	// Do a final run if you have any ops left
	if ops > 0 {
		_, err := mybulk.Run()
		checkError(err)
	}

	end := time.Now()
	log.Println(`End time config info Mongo: ` + end.Format(`1 January 2006, 15:04:05`))
	duration := time.Since(start)
	log.Print(`Time elapsed config info Mongo: `, duration.Minutes(), ` minutes`)

}

func UpsertConfigHistory(mongoSession *mgo.Session, bamiloCatalogConfigTable []bamilocatalogconfig.BamiloCatalogConfig, start time.Time) {

	bamiloCatalogConfigHistCollection := mongoSession.DB(`competition_analysis`).C(`bml_catalog_config_hist`)

	mybulk := bamiloCatalogConfigHistCollection.Bulk()
	mybulk.Unordered()
	ops := 0
	today := time.Now().Format(`01022006`)
	now := time.Now()

	for _, bamiloCatalogConfig := range bamiloCatalogConfigTable {
		iDMongoHist := bamiloCatalogConfig.IDBmlCatalogConfig + today
		if bamiloCatalogConfig.VisibleInShop == `1` {
			bamiloCatalogConfig.SetVisibleInShopTrue()
		}

		// only keep appropriate information for bamilocatalogconfig
		bamiloCatalogConfigHist := bamilocatalogconfig.BamiloCatalogConfig{
			IDMongoDb:          iDMongoHist,
			FKBmlCatalogConfig: bamiloCatalogConfig.IDBmlCatalogConfig,
			ConfigSnapshotAt:   now,

			VisibleInShopBool:  bamiloCatalogConfig.VisibleInShopBool,
			AvgPrice:           bamiloCatalogConfig.AvgPrice,
			AvgSpecialPrice:    bamiloCatalogConfig.AvgSpecialPrice,
			SumOfStockQuantity: bamiloCatalogConfig.SumOfStockQuantity,
			MinOfStockQuantity: bamiloCatalogConfig.MinOfStockQuantity,
		}

		var data interface{}
		data = bamiloCatalogConfigHist

		if ops > 950 {
			// We need to run our operations before we queue up too many
			_, err := mybulk.Run()
			checkError(err)

			// re-initialize the bulk
			mybulk = nil
			mybulk := bamiloCatalogConfigHistCollection.Bulk()
			mybulk.Unordered()

			// re-initialize ops
			ops = 0
		}

		ops++
		mybulk.Upsert(data)
	}
	// Do a final run if you have any ops left
	if ops > 0 {
		_, err := mybulk.Run()
		checkError(err)
	}

	end := time.Now()
	log.Println(`End time config info Mongo: ` + end.Format(`1 January 2006, 15:04:05`))
	duration := time.Since(start)
	log.Print(`Time elapsed config info Mongo: `, duration.Minutes(), ` minutes`)

}

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
