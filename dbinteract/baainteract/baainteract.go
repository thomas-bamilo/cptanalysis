package baainteract

import (
	"database/sql"
	"log"
	"time"

	"github.com/thomas-bamilo/commercial/competitionanalysis/bamilocatalogconfig"
)

// could speed up by selecting this beforehand for re-use: SELECT CONCAT(@p1, REPLACE(CONVERT (CHAR(10), GETDATE(), 101),'/',''));

func AddBamiloCatalogConfigTableToBaa(dbBaa *sql.DB, bamiloCatalogConfigTable []bamilocatalogconfig.BamiloCatalogConfig, start time.Time) {

	for _, bamiloCatalogConfig := range bamiloCatalogConfigTable {
		// if id_bml_catalog_config does not exist in cpt_bml_catalog_config
		// then create the new config in cpt_bml_catalog_config
		// and record the first historical data in cpt_bml_catalog_config_hist
		if !isConfigExist(dbBaa, bamiloCatalogConfig) {

			createConfigRecordHistory(dbBaa, bamiloCatalogConfig)

		} else {

			// else update the existing config in cpt_bml_catalog_config
			// and if historical data has already been recorded today
			// then update historical data in cpt_bml_catalog_config_hist
			// else record new historical data in cpt_bml_catalog_config_hist
			updateConfigRecordHistory(dbBaa, bamiloCatalogConfig)
		}
	}

	end := time.Now()
	log.Println(`End time SQL: ` + end.Format(`1 January 2006, 15:04:05`))
	duration := time.Since(start)
	log.Print(`Time elapsed SQL: `, duration.Minutes(), ` minutes`)
}

// check if the config already exists in cpt_bml_catalog_config
func isConfigExist(dbBaa *sql.DB, bamiloCatalogConfig bamilocatalogconfig.BamiloCatalogConfig) bool {

	var test string
	// store query in a string
	query := `
	SELECT 

	cbcc.id_bml_catalog_config

	FROM baa_application.commercial.cpt_bml_catalog_config cbcc
	WHERE cbcc.id_bml_catalog_config = @p1
	`

	err := dbBaa.QueryRow(query, bamiloCatalogConfig.IDBmlCatalogConfig).Scan(&test)
	// if result has no row, return false
	if err != nil {
		if err.Error() != `sql: no rows in result set` {
			log.Fatal(err.Error())
		} else {
			return false
		}

	}

	// else return true
	return true

}

// check if historical data has already been recorded today
func isConfigHistExist(dbBaa *sql.DB, bamiloCatalogConfig bamilocatalogconfig.BamiloCatalogConfig) bool {

	var test string
	// store query in a string
	query := `
	SELECT 

	cbchv.id_bml_catalog_config_hist

	FROM baa_application.commercial.cpt_bml_config_hist_view cbchv
	WHERE cbchv.id_bml_catalog_config_hist = CONCAT(@p1, REPLACE(CONVERT (CHAR(10), GETDATE(), 101),'/',''))
	`

	err := dbBaa.QueryRow(query, bamiloCatalogConfig.IDBmlCatalogConfig).Scan(&test)
	// if result has no row, return false
	if err != nil {
		if err.Error() != `sql: no rows in result set` {
			log.Fatal(err.Error())
		} else {
			return false
		}

	}

	// else return true
	return true

}

// create the new config in cpt_bml_catalog_config
// and record the first historical data in cpt_bml_catalog_config_hist
func createConfigRecordHistory(dbBaa *sql.DB, bamiloCatalogConfig bamilocatalogconfig.BamiloCatalogConfig) {

	// prepare statement
	createConfigRecordHistoryStr := `
		INSERT INTO baa_application.commercial.cpt_bml_catalog_config (
			id_bml_catalog_config
			,sku
			,sku_name
			,img_link
			,description 
			,short_description 
			,package_content 
			,product_warranty 
		  
			,bi_category_one_name 
			,bi_category_two_name 
			,bi_category_three_name 
		  
			,brand_name 
			,brand_name_en 
		  
			,supplier_name 
			,supplier_name_en) 
		VALUES (@p1,@p2,@p3,@p4,@p5,@p6,@p7,@p8,@p9,@p10,@p11,@p12,@p13,@p14,@p15);
		
		INSERT INTO baa_application.commercial.cpt_bml_catalog_config_hist (
			id_bml_catalog_config_hist
			,fk_bml_catalog_config

			,visible_in_shop

			,avg_price 
			,avg_special_price 
			,sum_of_stock_quantity 
			,min_of_stock_quantity) 
		VALUES (CONCAT(@p1, REPLACE(CONVERT (CHAR(10), GETDATE(), 101),'/','')),@p1,@p16,@p17,@p18,@p19,@p20);
		`
	createConfigRecordHistory, err := dbBaa.Prepare(createConfigRecordHistoryStr)
	defer createConfigRecordHistory.Close()
	checkError(err)

	// execute insert
	_, err = createConfigRecordHistory.Exec(
		bamiloCatalogConfig.IDBmlCatalogConfig,
		bamiloCatalogConfig.SKU,
		bamiloCatalogConfig.SKUName,
		bamiloCatalogConfig.ImgLink,
		bamiloCatalogConfig.Description,
		bamiloCatalogConfig.ShortDescription,
		bamiloCatalogConfig.PackageContent,
		bamiloCatalogConfig.ProductWarranty,

		bamiloCatalogConfig.BiCategoryOneName,
		bamiloCatalogConfig.BiCategoryTwoName,
		bamiloCatalogConfig.BiCategoryThreeName,

		bamiloCatalogConfig.BrandName,
		bamiloCatalogConfig.BrandNameEn,

		bamiloCatalogConfig.SupplierName,
		bamiloCatalogConfig.SupplierNameEn,

		bamiloCatalogConfig.VisibleInShop,

		bamiloCatalogConfig.AvgPrice,
		bamiloCatalogConfig.AvgSpecialPrice,
		bamiloCatalogConfig.SumOfStockQuantity,
		bamiloCatalogConfig.MinOfStockQuantity,
	)
	checkError(err)

}

// update the existing config in cpt_bml_catalog_config
// and if historical data has already been recorded today
// then update historical data in cpt_bml_catalog_config_hist
// else record new historical data in cpt_bml_catalog_config_hist
func updateConfigRecordHistory(dbBaa *sql.DB, bamiloCatalogConfig bamilocatalogconfig.BamiloCatalogConfig) {

	// prepare statement
	updateConfigRecordHistoryStr1 := `
	UPDATE baa_application.commercial.cpt_bml_catalog_config
				SET 
				config_updated_at = GETDATE()

				,sku_name = @p1
				,img_link = @p2
				,description  = @p3
				,short_description  = @p4
				,package_content  = @p5
				,product_warranty  = @p6
			  
				,bi_category_one_name  = @p7
				,bi_category_two_name  = @p8
				,bi_category_three_name  = @p9
			  
				,brand_name  = @p10
				,brand_name_en  = @p11
			  
				,supplier_name  = @p12
				,supplier_name_en = @p13
	
	WHERE id_bml_catalog_config = @p14;
		`
	updateConfigRecordHistory1, err := dbBaa.Prepare(updateConfigRecordHistoryStr1)
	defer updateConfigRecordHistory1.Close()
	checkError(err)

	// execute insert
	_, err = updateConfigRecordHistory1.Exec(
		bamiloCatalogConfig.SKUName,
		bamiloCatalogConfig.ImgLink,
		bamiloCatalogConfig.Description,
		bamiloCatalogConfig.ShortDescription,
		bamiloCatalogConfig.PackageContent,
		bamiloCatalogConfig.ProductWarranty,

		bamiloCatalogConfig.BiCategoryOneName,
		bamiloCatalogConfig.BiCategoryTwoName,
		bamiloCatalogConfig.BiCategoryThreeName,

		bamiloCatalogConfig.BrandName,
		bamiloCatalogConfig.BrandNameEn,

		bamiloCatalogConfig.SupplierName,
		bamiloCatalogConfig.SupplierNameEn,

		bamiloCatalogConfig.IDBmlCatalogConfig,
	)
	checkError(err)

	// if historical data has already been recorded today
	// then update historical data in cpt_bml_catalog_config_hist
	if isConfigHistExist(dbBaa, bamiloCatalogConfig) {
		// prepare statement
		updateConfigRecordHistoryStr2 := `
		UPDATE baa_application.commercial.cpt_bml_catalog_config_hist
				SET 
					visible_in_shop = @p2
		
					,avg_price = @p3
					,avg_special_price = @p4
					,sum_of_stock_quantity = @p5
					,min_of_stock_quantity = @p6
		
		WHERE id_bml_catalog_config_hist = CONCAT(@p1, REPLACE(CONVERT (CHAR(10), GETDATE(), 101),'/',''));
		`
		updateConfigRecordHistory2, err := dbBaa.Prepare(updateConfigRecordHistoryStr2)
		defer updateConfigRecordHistory2.Close()
		checkError(err)

		// execute insert
		_, err = updateConfigRecordHistory2.Exec(
			bamiloCatalogConfig.IDBmlCatalogConfig,

			bamiloCatalogConfig.VisibleInShop,

			bamiloCatalogConfig.AvgPrice,
			bamiloCatalogConfig.AvgSpecialPrice,
			bamiloCatalogConfig.SumOfStockQuantity,
			bamiloCatalogConfig.MinOfStockQuantity,
		)
		checkError(err)

		// else record new historical data in cpt_bml_catalog_config_hist
	} else {
		// prepare statement
		updateConfigRecordHistoryStr3 := `
		INSERT INTO baa_application.commercial.cpt_bml_catalog_config_hist (
			
			id_bml_catalog_config_hist
			,fk_bml_catalog_config

			,visible_in_shop

			,avg_price 
			,avg_special_price 
			,sum_of_stock_quantity 
			,min_of_stock_quantity) 
		VALUES (CONCAT(@p1, REPLACE(CONVERT (CHAR(10), GETDATE(), 101),'/','')),@p1,@p2,@p3,@p4,@p5,@p6);
`
		updateConfigRecordHistory3, err := dbBaa.Prepare(updateConfigRecordHistoryStr3)
		defer updateConfigRecordHistory3.Close()
		checkError(err)

		// execute insert
		_, err = updateConfigRecordHistory3.Exec(
			bamiloCatalogConfig.IDBmlCatalogConfig,

			bamiloCatalogConfig.VisibleInShop,

			bamiloCatalogConfig.AvgPrice,
			bamiloCatalogConfig.AvgSpecialPrice,
			bamiloCatalogConfig.SumOfStockQuantity,
			bamiloCatalogConfig.MinOfStockQuantity,

			bamiloCatalogConfig.CountOfSoi,
			bamiloCatalogConfig.SumOfUnitPrice,
			bamiloCatalogConfig.SumOfPaidPrice,
			bamiloCatalogConfig.SumOfCouponMoneyValue,
			bamiloCatalogConfig.SumOfCartRuleDiscount,
		)
		checkError(err)
	}
}

// -----------------------------------------------------------------------------------------------------------------------------------------------------------------------------
// Update sales historical data

func UpdateBamiloCatalogConfigSales(dbBaa *sql.DB, bamiloCatalogConfigSalesTable []bamilocatalogconfig.BamiloCatalogConfig) {

	for _, bamiloCatalogConfigSales := range bamiloCatalogConfigSalesTable {
		// if config historical data exists (and it always should)
		// then update historical sales data in cpt_bml_catalog_config_hist
		if isConfigHistExist(dbBaa, bamiloCatalogConfigSales) {

			log.Println(`Updated`, bamiloCatalogConfigSales.IDBmlCatalogConfig, bamiloCatalogConfigSales.ConfigSnapshotAt)
			updateConfigSalesHistory(dbBaa, bamiloCatalogConfigSales)

		} else {

			log.Println(`Not updated`, bamiloCatalogConfigSales.IDBmlCatalogConfig, bamiloCatalogConfigSales.ConfigSnapshotAt)

		}
	}

}

func updateConfigSalesHistory(dbBaa *sql.DB, bamiloCatalogConfig bamilocatalogconfig.BamiloCatalogConfig) {

	// prepare statement
	updateConfigSalesHistoryStr := `
	UPDATE baa_application.commercial.cpt_bml_catalog_config_hist
			SET 	
				count_of_soi = @p3
				,sum_of_paid_price = @p4
				,sum_of_unit_price = @p5
				,sum_of_coupon_money_value = @p6
				,sum_of_cart_rule_discount = @p7
	
	WHERE id_bml_catalog_config_hist = CONCAT(@p1, REPLACE(CONVERT (CHAR(10), GETDATE(), 101),'/',''));
	`
	updateConfigSalesHistory, err := dbBaa.Prepare(updateConfigSalesHistoryStr)
	defer updateConfigSalesHistory.Close()
	checkError(err)

	// execute update
	_, err = updateConfigSalesHistory.Exec(
		bamiloCatalogConfig.IDBmlCatalogConfig,
		bamiloCatalogConfig.ConfigSnapshotAt,

		bamiloCatalogConfig.CountOfSoi,
		bamiloCatalogConfig.SumOfPaidPrice,
		bamiloCatalogConfig.SumOfUnitPrice,
		bamiloCatalogConfig.SumOfCouponMoneyValue,
		bamiloCatalogConfig.SumOfCartRuleDiscount,
	)
	checkError(err)

}

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
