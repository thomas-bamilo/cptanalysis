package bobinteract

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"

	"github.com/thomas-bamilo/commercial/competitionanalysis/bamilocatalogconfig"
	"github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/baainteract"
	"github.com/thomas-bamilo/commercial/competitionanalysis/dbinteract/elasticinteract"
	elastic "gopkg.in/olivere/elastic.v5"
)

func BamiloCatalogConfigFromBobToElasticAndBaa(dbBob *sql.DB, elasticClient *elastic.Client, dbBaa *sql.DB, ctx context.Context) {

	stmt, err := dbBob.Prepare(`
		SELECT 

		cc.id_catalog_config
		,cc.sku
		,COALESCE(cc.name,"") sku_name
		,COALESCE(CONCAT('https://media.bamilo.com/p/',cb.url_key,'-',RIGHT(UNIX_TIMESTAMP(cpi.updated_at),4),'-',REVERSE(cs.fk_catalog_config),'-1-product.jpg'),"") img_link
		,COALESCE(cc.description,"") description
		,COALESCE(cc.short_description,"") short_description
		,COALESCE(cc.package_content,"") package_content
		,COALESCE(cc.product_warranty,"") product_warranty

		,COALESCE(AVG(cs.price),"") avg_price
		,COALESCE(AVG(cs.special_price),"") avg_special_price
		,COALESCE(SUM(cs2.quantity),"") sum_of_stock_quantity
    	,COALESCE(MIN(cs2.quantity),"") min_of_stock_quantity
		
		,COALESCE(bi_one.name,"") bi_category_one_name
		,COALESCE(bi_two.name,"") bi_category_two_name
		,COALESCE(bi_three.name,"") bi_category_three_name
		
		,COALESCE(cb.name,"") brand_name
		,COALESCE(cb.name_en,"") brand_name_en
		
		,COALESCE(s.name,"") supplier_name
		,COALESCE(s.name_en,"") supplier_name_en
		  
		FROM catalog_simple cs
		JOIN catalog_config cc
		ON cs.fk_catalog_config = cc.id_catalog_config
		JOIN catalog_category_bi bi_one
		ON cc.bi_category_one = bi_one.id_catalog_category_bi
		JOIN catalog_category_bi bi_two
		ON cc.bi_category_two = bi_two.id_catalog_category_bi
		JOIN catalog_category_bi bi_three
		ON cc.bi_category_three = bi_three.id_catalog_category_bi
		JOIN catalog_brand cb
		ON cc.fk_catalog_brand = cb.id_catalog_brand
		JOIN catalog_source cs1
		ON cs.id_catalog_simple = cs1.fk_catalog_simple
		JOIN supplier s
		ON cs1.fk_supplier = s.id_supplier
		JOIN catalog_stock cs2
		ON cs1.id_catalog_source = cs2.fk_catalog_source
		JOIN catalog_product_image cpi
		ON cpi.fk_catalog_config = cc.id_catalog_config
		
		GROUP BY cc.id_catalog_config;
	 `)
	checkError(err)
	defer stmt.Close()

	rows, err := stmt.Query()
	checkError(err)
	defer rows.Close()

	bamiloCatalogConfig := bamilocatalogconfig.BamiloCatalogConfig{}

	for rows.Next() {
		// SKUHistoryUpdatedAt is saved here but will only be saved to Elastic if the SKU is actually updated
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
			// quantitative information
			&bamiloCatalogConfig.AvgPrice,
			&bamiloCatalogConfig.AvgSpecialPrice,
			&bamiloCatalogConfig.SumOfStockQuantity,
			&bamiloCatalogConfig.MinOfStockQuantity,
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
		)
		checkError(err)

		// only keep appropriate information for elastic
		bamiloCatalogConfigElastic := bamilocatalogconfig.BamiloCatalogConfig{
			IDBmlCatalogConfig: bamiloCatalogConfig.IDBmlCatalogConfig,
			SKUName:            bamiloCatalogConfig.SKUName,
			Description:        bamiloCatalogConfig.Description,
			ShortDescription:   bamiloCatalogConfig.ShortDescription,
			PackageContent:     bamiloCatalogConfig.PackageContent,
			ProductWarranty:    bamiloCatalogConfig.ProductWarranty,

			AvgPrice:        bamiloCatalogConfig.AvgPrice,
			AvgSpecialPrice: bamiloCatalogConfig.AvgSpecialPrice,

			BiCategoryOneName:   bamiloCatalogConfig.BiCategoryOneName,
			BiCategoryTwoName:   bamiloCatalogConfig.BiCategoryTwoName,
			BiCategoryThreeName: bamiloCatalogConfig.BiCategoryThreeName,

			BrandName:   bamiloCatalogConfig.BrandName,
			BrandNameEn: bamiloCatalogConfig.BrandNameEn,

			SupplierName:   bamiloCatalogConfig.SupplierName,
			SupplierNameEn: bamiloCatalogConfig.SupplierNameEn,
		}

		// convert information to JSON
		bamiloCatalogConfigElasticByte, err := json.Marshal(bamiloCatalogConfigElastic)
		checkError(err)
		bamiloCatalogConfigElasticJSON := string(bamiloCatalogConfigElasticByte)
		// add information to elastic
		go elasticinteract.AddBamiloCatalogConfigToIndex(elasticClient, ctx, bamiloCatalogConfigElasticJSON)

		baainteract.AddBamiloCatalogConfigToBaa(dbBaa, bamiloCatalogConfig)

	}

}

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
