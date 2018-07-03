package bobinteract

import (
	"database/sql"
	"log"
	"time"

	"github.com/thomas-bamilo/commercial/competitionanalysis/bamilocatalogconfig"
)

func GetBamiloCatalogConfigTable(dbBob *sql.DB) []bamilocatalogconfig.BamiloCatalogConfig {

	stmt, err := dbBob.Prepare(`
		-- it is fine if several SKUs have the same name, keep structure similar to BOB - plus it will be much clearer for business users
  	SELECT
-- qualitative information
		cc.id_catalog_config
		,cc.sku
		,COALESCE(cc.name,'') sku_name
		,COALESCE(CONCAT('https://media.bamilo.com/p/',cb.url_key,'-',RIGHT(UNIX_TIMESTAMP(cpi.updated_at),4),'-',REVERSE(cs.fk_catalog_config),'-1-product.jpg'),'') img_link
		,COALESCE(cc.description,'') description
		,COALESCE(cc.short_description,'') short_description
		,COALESCE(cc.package_content,'') package_content
		,COALESCE(cc.product_warranty,'') product_warranty
-- category
		,COALESCE(bi_one.name,'') bi_category_one_name
		,COALESCE(bi_two.name,'') bi_category_two_name
		,COALESCE(bi_three.name,'') bi_category_three_name
-- brand
		,COALESCE(cb.name,'') brand_name
		,COALESCE(cb.name_en,'') brand_name_en
-- supplier
		,COALESCE(s.name,'') supplier_name
		,COALESCE(s.name_en,'') supplier_name_en
-- historical data
    ,COALESCE(vccv.visible_in_shop,0)
		,COALESCE(FLOOR(AVG(cs.price)),0) avg_price
		,COALESCE(FLOOR(AVG(cs.special_price)),0) avg_special_price
		,COALESCE(FLOOR(SUM(cs2.quantity)),0) sum_of_stock_quantity
    	,COALESCE(FLOOR(MIN(cs2.quantity)),0) min_of_stock_quantity
  
		FROM catalog_simple cs
		JOIN catalog_config cc
		ON cs.fk_catalog_config = cc.id_catalog_config
		JOIN catalog_brand cb
		ON cc.fk_catalog_brand = cb.id_catalog_brand
		JOIN catalog_product_image cpi
		ON cpi.fk_catalog_config = cc.id_catalog_config
    JOIN bob_live_ir.view_catalog_config_visibility vccv
		ON vccv.sku = cc.sku
    LEFT JOIN catalog_category_bi bi_one
		ON cc.bi_category_one = bi_one.id_catalog_category_bi
		LEFT JOIN catalog_category_bi bi_two
		ON cc.bi_category_two = bi_two.id_catalog_category_bi
		LEFT JOIN catalog_category_bi bi_three
		ON cc.bi_category_three = bi_three.id_catalog_category_bi -- some products do not have bi_category_three, just to make sure, also left join other bi_category
		LEFT JOIN catalog_source cs1
		ON cs.id_catalog_simple = cs1.fk_catalog_simple
		LEFT JOIN catalog_stock cs2
		ON cs1.id_catalog_source = cs2.fk_catalog_source
    LEFT JOIN supplier s
		ON cs1.fk_supplier = s.id_supplier
    LEFT JOIN (
    SELECT 
     soi.sku
     ,soi.created_at created_at -- not possible to refresh sales here, should be done with sales script
    FROM sales_order_item soi
    WHERE soi.created_at >= NOW()-INTERVAL 30 DAY
    GROUP BY soi.sku, soi.created_at
    ) sales_order_sku
    ON sales_order_sku.sku = cs.sku
   
    WHERE (sales_order_sku.created_at >= NOW()-INTERVAL 30 DAY
    OR vccv.visible_in_shop = 1)

		GROUP BY cc.id_catalog_config;
	 `)
	checkError(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	checkError(err)
	defer rows.Close()

	var bamiloCatalogConfigTable []bamilocatalogconfig.BamiloCatalogConfig
	bamiloCatalogConfig := bamilocatalogconfig.BamiloCatalogConfig{}

	for rows.Next() {

		err := rows.Scan(
			// qualitative information
			&bamiloCatalogConfig.IDBmlCatalogConfig,
			&bamiloCatalogConfig.SKU,
			&bamiloCatalogConfig.SKUName,
			&bamiloCatalogConfig.ImgLink,
			&bamiloCatalogConfig.Description,
			&bamiloCatalogConfig.ShortDescription,
			&bamiloCatalogConfig.PackageContent,
			&bamiloCatalogConfig.ProductWarranty,
			// category
			&bamiloCatalogConfig.BiCategoryOneName,
			&bamiloCatalogConfig.BiCategoryTwoName,
			&bamiloCatalogConfig.BiCategoryThreeName,
			// brand
			&bamiloCatalogConfig.BrandName,
			&bamiloCatalogConfig.BrandNameEn,
			// supplier
			&bamiloCatalogConfig.SupplierName,
			&bamiloCatalogConfig.SupplierNameEn,
			// historical data
			&bamiloCatalogConfig.VisibleInShop,
			&bamiloCatalogConfig.AvgPrice,
			&bamiloCatalogConfig.AvgSpecialPrice,
			&bamiloCatalogConfig.SumOfStockQuantity,
			&bamiloCatalogConfig.MinOfStockQuantity,
		)
		checkError(err)

		bamiloCatalogConfigTable = append(bamiloCatalogConfigTable, bamiloCatalogConfig)

	}

	return bamiloCatalogConfigTable

}

// GetBamiloCatalogConfigSalesTable retrieves the sum of sales data across a certain period for each SKU from BOB
func GetBamiloCatalogConfigSalesTable(dbBob *sql.DB) []bamilocatalogconfig.BamiloCatalogConfig {

	stmt, err := dbBob.Prepare(`
		SELECT

		cc.id_catalog_config

		,COUNT(DISTINCT soi.id_sales_order_item) count_of_soi
		,COALESCE(FLOOR(SUM(soi.paid_price))) sum_of_paid_price
		,COALESCE(FLOOR(SUM(soi.unit_price))) sum_of_unit_price
		,COALESCE(FLOOR(SUM(soi.coupon_money_value))) sum_of_coupon_money_value
		,COALESCE(FLOOR(SUM(soi.cart_rule_discount))) sum_of_cart_rule_discount
	
		FROM sales_order_item soi
		JOIN catalog_simple cs
		ON soi.sku = cs.sku
		JOIN catalog_config cc
		ON cs.fk_catalog_config = cc.id_catalog_config
	
		WHERE CAST(soi.created_at AS DATE) = CAST(NOW()-INTERVAL 1 DAY AS DATE)
	
		GROUP BY cc.id_catalog_config;
	 `)
	checkError(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	checkError(err)
	defer rows.Close()

	var bamiloCatalogConfigSalesTable []bamilocatalogconfig.BamiloCatalogConfig
	bamiloCatalogConfigSales := bamilocatalogconfig.BamiloCatalogConfig{}

	for rows.Next() {

		err := rows.Scan(
			&bamiloCatalogConfigSales.IDBmlCatalogConfig,

			&bamiloCatalogConfigSales.CountOfSoi,
			&bamiloCatalogConfigSales.SumOfPaidPrice,
			&bamiloCatalogConfigSales.SumOfUnitPrice,
			&bamiloCatalogConfigSales.SumOfCouponMoneyValue,
			&bamiloCatalogConfigSales.SumOfCartRuleDiscount,
		)
		checkError(err)

		bamiloCatalogConfigSalesTable = append(bamiloCatalogConfigSalesTable, bamiloCatalogConfigSales)

	}

	//log.Println(` Length of the table: ` + strconv.Itoa(len(bamiloCatalogConfigSalesTable)))

	return bamiloCatalogConfigSalesTable

}

// GetBamiloCatalogConfigSalesHistTable retrieves the sum of sales for each day and each SKU from BOB
func GetBamiloCatalogConfigSalesHistTable(dbBob *sql.DB) []bamilocatalogconfig.BamiloCatalogConfig {

	stmt, err := dbBob.Prepare(`
		SELECT

		cc.id_catalog_config
		,DATE_FORMAT(soi.created_at, "%m/%d/%Y") config_snapshot_at

		,COUNT(DISTINCT soi.id_sales_order_item) count_of_soi
		,COALESCE(FLOOR(SUM(soi.paid_price))) sum_of_paid_price
		,COALESCE(FLOOR(SUM(soi.unit_price))) sum_of_unit_price
		,COALESCE(FLOOR(SUM(soi.coupon_money_value))) sum_of_coupon_money_value
		,COALESCE(FLOOR(SUM(soi.cart_rule_discount))) sum_of_cart_rule_discount
	
		FROM sales_order_item soi
		JOIN catalog_simple cs
		ON soi.sku = cs.sku
		JOIN catalog_config cc
		ON cs.fk_catalog_config = cc.id_catalog_config
	
		WHERE CAST(soi.created_at AS DATE) >= CAST(NOW()-INTERVAL 3 DAY AS DATE)
	
		GROUP BY cc.id_catalog_config, CAST(soi.created_at AS DATE);
	 `)
	checkError(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	checkError(err)
	defer rows.Close()

	var bamiloCatalogConfigSalesHistTable []bamilocatalogconfig.BamiloCatalogConfig
	var configSnapshotAtStr string
	bamiloCatalogConfigSalesHist := bamilocatalogconfig.BamiloCatalogConfig{}

	for rows.Next() {

		err := rows.Scan(
			&bamiloCatalogConfigSalesHist.IDBmlCatalogConfig,
			&configSnapshotAtStr,

			&bamiloCatalogConfigSalesHist.CountOfSoi,
			&bamiloCatalogConfigSalesHist.SumOfPaidPrice,
			&bamiloCatalogConfigSalesHist.SumOfUnitPrice,
			&bamiloCatalogConfigSalesHist.SumOfCouponMoneyValue,
			&bamiloCatalogConfigSalesHist.SumOfCartRuleDiscount,
		)
		checkError(err)

		bamiloCatalogConfigSalesHist.ConfigSnapshotAt, err = time.Parse(`01/02/2006`, configSnapshotAtStr)
		checkError(err)

		bamiloCatalogConfigSalesHistTable = append(bamiloCatalogConfigSalesHistTable, bamiloCatalogConfigSalesHist)

	}

	return bamiloCatalogConfigSalesHistTable

}

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
