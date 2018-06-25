package baainteract

import (
	"database/sql"
	"log"

	"github.com/thomas-bamilo/commercial/competitionanalysis/bamilocatalogconfig"
)

func AddBamiloCatalogConfigToBaa(dbBaa *sql.DB, bamiloCatalogConfig bamilocatalogconfig.BamiloCatalogConfig) {

	// if id_bml_catalog_config does not exist in cpt_bml_catalog_config
	// then create the new config in cpt_bml_catalog_config
	// and record the first historical data in cpt_bml_catalog_config_hist
	if isConfigExist(dbBaa, bamiloCatalogConfig) {

		createConfigRecordHistory(dbBaa, bamiloCatalogConfig)

	} else {

		// else update the existing config in cpt_bml_catalog_config
		// and if historical data has already been recorded today
		// then update historical data in cpt_bml_catalog_config_hist
		// else record new historical data in cpt_bml_catalog_config_hist
		updateConfigRecordHistory(dbBaa, bamiloCatalogConfig)
	}
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

	cbcch.fk_bml_catalog_config

	FROM baa_application.commercial.cpt_bml_catalog_config_hist cbcch
	WHERE cbcch.fk_bml_catalog_config = @p1
	AND cbcch.config_snapshot_at = CONVERT(DATE,GETDATE())
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
			fk_bml_catalog_config

			,visible_in_shop

			,avg_price 
			,avg_special_price 
			,sum_of_stock_quantity 
			,min_of_stock_quantity 
		  
			,count_of_soi 
			,sum_of_unit_price 
			,sum_of_paid_price 
			,sum_of_coupon_money_value 
			,sum_of_cart_rule_discount ) 
		VALUES (@p1,@p16,@p17,@p18,@p19,@p20,@p21,@p22,@p23,@p24,@p25);
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

		bamiloCatalogConfig.CountOfSoi,
		bamiloCatalogConfig.SumOfUnitPrice,
		bamiloCatalogConfig.SumOfPaidPrice,
		bamiloCatalogConfig.SumOfCouponMoneyValue,
		bamiloCatalogConfig.SumOfCartRuleDiscount,
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
				sku_name = @p1
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
	
	WHERE id_bml_catalog_config = @p14);
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
				  
					,count_of_soi = @p7
					,sum_of_unit_price = @p8
					,sum_of_paid_price = @p9
					,sum_of_coupon_money_value = @p10
					,sum_of_cart_rule_discount= @p11
		
		WHERE fk_bml_catalog_config = @p1
		AND config_snapshot_at = CONVERT(DATE,GETDATE());
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

			bamiloCatalogConfig.CountOfSoi,
			bamiloCatalogConfig.SumOfUnitPrice,
			bamiloCatalogConfig.SumOfPaidPrice,
			bamiloCatalogConfig.SumOfCouponMoneyValue,
			bamiloCatalogConfig.SumOfCartRuleDiscount,
		)
		checkError(err)

		// else record new historical data in cpt_bml_catalog_config_hist
	} else {
		// prepare statement
		updateConfigRecordHistoryStr3 := `
		INSERT INTO baa_application.commercial.cpt_bml_catalog_config_hist (
			fk_bml_catalog_config

			,visible_in_shop

			,avg_price 
			,avg_special_price 
			,sum_of_stock_quantity 
			,min_of_stock_quantity 
		
			,count_of_soi 
			,sum_of_unit_price 
			,sum_of_paid_price 
			,sum_of_coupon_money_value 
			,sum_of_cart_rule_discount ) 
		VALUES (@p1,@p2,@p3,@p4,@p5,@p6,@p7,@p8,@p9,@p10,@p11);
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

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
