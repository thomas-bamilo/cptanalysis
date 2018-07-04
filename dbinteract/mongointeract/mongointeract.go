package mongointeract

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/globalsign/mgo"
	"github.com/thomas-bamilo/commercial/competitionanalysis/bamilocatalogconfig"
	"github.com/thomas-bamilo/nosql/mongobulk"
)

func UpsertConfigInfo(mongoSession *mgo.Session, bamiloCatalogConfigTable []bamilocatalogconfig.BamiloCatalogConfig, start time.Time, wg *sync.WaitGroup) {

	defer wg.Done()

	now := time.Now()

	for _, bamiloCatalogConfig := range bamiloCatalogConfigTable {

		// only keep appropriate information for bamilocatalogconfig
		bamiloCatalogConfigInfo := bamilocatalogconfig.BamiloCatalogConfigInfo{
			IDBmlCatalogConfig: bamiloCatalogConfig.IDBmlCatalogConfig,
			ConfigSnapshotAt:   now,
			SKUName:            bamiloCatalogConfig.SKUName,
			ImgLink:            bamiloCatalogConfig.ImgLink,
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

		if bamiloCatalogConfig.VisibleInShop == `1` {
			bamiloCatalogConfigInfo.SetVisibleInShopTrue()
		}

		bamiloCatalogConfigInfo.UpsertConfigInfo(mongoSession)

	}

	end := time.Now()
	log.Println(`End time config info Mongo: ` + end.Format(`1 January 2006, 15:04:05`))
	duration := time.Since(start)
	log.Print(`Time elapsed config info Mongo: `, duration.Minutes(), ` minutes`)

}

func UpsertConfigInfoHist(mongoSession *mgo.Session, bamiloCatalogConfigTable []bamilocatalogconfig.BamiloCatalogConfig, start time.Time, wg *sync.WaitGroup) {

	defer wg.Done()

	now := time.Now()

	for _, bamiloCatalogConfig := range bamiloCatalogConfigTable {

		iDBmlCatalogConfigHist, err := strconv.Atoi(strconv.Itoa(bamiloCatalogConfig.IDBmlCatalogConfig) + now.Format(`01022006`))
		checkError(err)

		// only keep appropriate information for bamilocatalogconfig
		bamiloCatalogConfigInfoHist := bamilocatalogconfig.BamiloCatalogConfigInfoHist{
			IDBmlCatalogConfigHist: iDBmlCatalogConfigHist,
			FKBmlCatalogConfig:     bamiloCatalogConfig.IDBmlCatalogConfig,
			ConfigSnapshotAt:       now,

			VisibleInShopBool:  bamiloCatalogConfig.VisibleInShopBool,
			AvgPrice:           bamiloCatalogConfig.AvgPrice,
			AvgSpecialPrice:    bamiloCatalogConfig.AvgSpecialPrice,
			SumOfStockQuantity: bamiloCatalogConfig.SumOfStockQuantity,
			MinOfStockQuantity: bamiloCatalogConfig.MinOfStockQuantity,
		}

		if bamiloCatalogConfig.VisibleInShop == `1` {
			bamiloCatalogConfigInfoHist.SetVisibleInShopTrue()
		}

		bamiloCatalogConfigInfoHist.UpsertConfigInfoHist(mongoSession)

	}

	end := time.Now()
	log.Println(`End time config hist Mongo: ` + end.Format(`1 January 2006, 15:04:05`))
	duration := time.Since(start)
	log.Print(`Time elapsed config hist Mongo: `, duration.Minutes(), ` minutes`)

}

// sales ------------------------------------------------------------------------------

func UpsertConfigSales(mongoSession *mgo.Session, bamiloCatalogConfigSalesTable []bamilocatalogconfig.BamiloCatalogConfig, start time.Time, wg *sync.WaitGroup) {

	defer wg.Done()

	now := time.Now()

	c := mongoSession.DB("competition_analysis").C("bml_catalog_config")

	config := mongobulk.Config{OpsPerBatch: 950}

	mongoBulk := mongobulk.New(c, config)

	for _, bamiloCatalogConfigSales := range bamiloCatalogConfigSalesTable {

		// only keep appropriate information for bamilocatalogconfig
		bamiloCatalogConfigSales := bamilocatalogconfig.BamiloCatalogConfigSales{
			IDBmlCatalogConfig: bamiloCatalogConfigSales.IDBmlCatalogConfig,
			ConfigSnapshotAt:   now,

			CountOfSoi:            bamiloCatalogConfigSales.CountOfSoi,
			SumOfUnitPrice:        bamiloCatalogConfigSales.SumOfUnitPrice,
			SumOfPaidPrice:        bamiloCatalogConfigSales.SumOfPaidPrice,
			SumOfCouponMoneyValue: bamiloCatalogConfigSales.SumOfCouponMoneyValue,
			SumOfCartRuleDiscount: bamiloCatalogConfigSales.SumOfCartRuleDiscount,
		}

		bamiloCatalogConfigSales.UpsertConfigSales(mongoBulk)

	}

	end := time.Now()
	log.Println(`End time config sales Mongo: ` + end.Format(`1 January 2006, 15:04:05`))
	duration := time.Since(start)
	log.Print(`Time elapsed config sales Mongo: `, duration.Minutes(), ` minutes`)

	err := mongoBulk.Finish()
	checkError(err)

}

func UpsertConfigSalesHist(mongoSession *mgo.Session, bamiloCatalogConfigSalesHistTable []bamilocatalogconfig.BamiloCatalogConfig, start time.Time, wg *sync.WaitGroup) {

	defer wg.Done()

	c := mongoSession.DB("competition_analysis").C("bml_catalog_config_hist")

	config := mongobulk.Config{OpsPerBatch: 950}

	mongoBulk := mongobulk.New(c, config)

	for _, bamiloCatalogConfigSalesHist := range bamiloCatalogConfigSalesHistTable {

		iDBmlCatalogConfigHist, err := strconv.Atoi(strconv.Itoa(bamiloCatalogConfigSalesHist.IDBmlCatalogConfig) + bamiloCatalogConfigSalesHist.ConfigSnapshotAt.Format(`01022006`))
		checkError(err)
		// only keep appropriate information for bamilocatalogconfig
		bamiloCatalogConfigSalesHist := bamilocatalogconfig.BamiloCatalogConfigSalesHist{
			IDBmlCatalogConfigHist: iDBmlCatalogConfigHist,
			FKBmlCatalogConfig:     bamiloCatalogConfigSalesHist.IDBmlCatalogConfig,
			ConfigSnapshotAt:       bamiloCatalogConfigSalesHist.ConfigSnapshotAt,

			CountOfSoi:            bamiloCatalogConfigSalesHist.CountOfSoi,
			SumOfUnitPrice:        bamiloCatalogConfigSalesHist.SumOfUnitPrice,
			SumOfPaidPrice:        bamiloCatalogConfigSalesHist.SumOfPaidPrice,
			SumOfCouponMoneyValue: bamiloCatalogConfigSalesHist.SumOfCouponMoneyValue,
			SumOfCartRuleDiscount: bamiloCatalogConfigSalesHist.SumOfCartRuleDiscount,
		}

		bamiloCatalogConfigSalesHist.UpsertConfigSalesHist(mongoBulk)

	}
	end := time.Now()
	log.Println(`End time config sales Mongo: ` + end.Format(`1 January 2006, 15:04:05`))
	duration := time.Since(start)
	log.Print(`Time elapsed config sales Mongo: `, duration.Minutes(), ` minutes`)

	err := mongoBulk.Finish()
	checkError(err)

}

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
