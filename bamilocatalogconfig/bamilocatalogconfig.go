package bamilocatalogconfig

import (
	"log"
	"time"

	"github.com/globalsign/mgo"
)

type BamiloCatalogConfig struct {
	// qualitative information
	IDBmlCatalogConfig     int       `json:"id_bml_catalog_config" bson:"fk_bml_catalog_config"`
	FKBmlCatalogConfig     int       `json:"fk_bml_catalog_config" bson:"fk_bml_catalog_config"`
	IDBmlCatalogConfigHist int       `json:"id_bml_catalog_config_hist" bson:"id_bml_catalog_config_hist"`
	ConfigSnapshotAt       time.Time `json:"config_snapshot_at" bson:"config_snapshot_at"` // for sales: order_created_at
	SKU                    string    `json:"sku" bson:"sku"`
	SKUName                string    `json:"sku_name" bson:"sku_name"`
	ImgLink                string    `json:"img_link" bson:"img_link"`
	Description            string    `json:"description" bson:"description"`
	ShortDescription       string    `json:"short_description" bson:"short_description"`
	PackageContent         string    `json:"package_content" bson:"package_content"`
	ProductWarranty        string    `json:"product_warranty" bson:"product_warranty"`
	// category
	BiCategoryOneName   string `json:"bi_category_one_name" bson:"bi_category_one_name"`
	BiCategoryTwoName   string `json:"bi_category_two_name" bson:"bi_category_two_name"`
	BiCategoryThreeName string `json:"bi_category_three_name" bson:"bi_category_three_name"`
	// brand
	BrandName   string `json:"brand_name" bson:"brand_name"`
	BrandNameEn string `json:"brand_name_en" bson:"brand_name_en"`
	// supplier
	SupplierName   string `json:"supplier_name" bson:"supplier_name"`
	SupplierNameEn string `json:"supplier_name_en" bson:"supplier_name_en"`
	// historical visibility
	VisibleInShop     string
	VisibleInShopBool bool `json:"visible_in_shop" bson:"visible_in_shop"`
	// historical price and quantity
	AvgPrice           int `json:"avg_price" bson:"avg_price"`
	AvgSpecialPrice    int `json:"avg_special_price" bson:"avg_special_price"`
	SumOfStockQuantity int `json:"sum_of_stock_quantity" bson:"sum_of_stock_quantity"`
	MinOfStockQuantity int `json:"min_of_stock_quantity" bson:"min_of_stock_quantity"`
	// sales
	CountOfSoi            int `json:"count_of_soi" bson:"count_of_soi"`
	SumOfUnitPrice        int `json:"sum_of_unit_price" bson:"sum_of_unit_price"`
	SumOfPaidPrice        int `json:"sum_of_paid_price" bson:"sum_of_paid_price"`
	SumOfCouponMoneyValue int `json:"sum_of_coupon_money_value" bson:"sum_of_coupon_money_value"`
	SumOfCartRuleDiscount int `json:"sum_of_cart_rule_discount" bson:"sum_of_cart_rule_discount"`
}

// elastic -----------------------------------------------------------------------------------------------------------------------------------------------

type BamiloCatalogConfigElastic struct {
	// qualitative information
	SKUName          string `json:"sku_name"`
	ImgLink          string `json:"img_link"`
	Description      string `json:"description"`
	ShortDescription string `json:"short_description"`
	PackageContent   string `json:"package_content"`
	ProductWarranty  string `json:"product_warranty"`
	// category
	BiCategoryOneName   string `json:"bi_category_one_name"`
	BiCategoryTwoName   string `json:"bi_category_two_name"`
	BiCategoryThreeName string `json:"bi_category_three_name"`
	// brand
	BrandName   string `json:"brand_name"`
	BrandNameEn string `json:"brand_name_en"`
	// supplier
	SupplierName   string `json:"supplier_name"`
	SupplierNameEn string `json:"supplier_name_en"`
	// historical price and quantity
	AvgPrice        int `json:"avg_price"`
	AvgSpecialPrice int `json:"avg_special_price"`
}

// mongo ----------------------------------------------------------------------------------------------------------------------------------------------

// config info -------------------------------------------------------------

type BamiloCatalogConfigInfo struct {
	// qualitative information
	IDBmlCatalogConfig int       `json:"id_bml_catalog_config" bson:"id_bml_catalog_config"`
	ConfigSnapshotAt   time.Time `json:"config_snapshot_at" bson:"config_snapshot_at"` // for sales: order_created_at
	SKUName            string    `json:"sku_name" bson:"sku_name"`
	ImgLink            string    `json:"img_link" bson:"img_link"`
	Description        string    `json:"description" bson:"description"`
	ShortDescription   string    `json:"short_description" bson:"short_description"`
	PackageContent     string    `json:"package_content" bson:"package_content"`
	ProductWarranty    string    `json:"product_warranty" bson:"product_warranty"`
	// category
	BiCategoryOneName   string `json:"bi_category_one_name" bson:"bi_category_one_name"`
	BiCategoryTwoName   string `json:"bi_category_two_name" bson:"bi_category_two_name"`
	BiCategoryThreeName string `json:"bi_category_three_name" bson:"bi_category_three_name"`
	// brand
	BrandName   string `json:"brand_name" bson:"brand_name"`
	BrandNameEn string `json:"brand_name_en" bson:"brand_name_en"`
	// supplier
	SupplierName   string `json:"supplier_name" bson:"supplier_name"`
	SupplierNameEn string `json:"supplier_name_en" bson:"supplier_name_en"`
	// historical visibility
	VisibleInShopBool bool `json:"visible_in_shop" bson:"visible_in_shop"`
	// historical price and quantity
	AvgPrice           int `json:"avg_price" bson:"avg_price"`
	AvgSpecialPrice    int `json:"avg_special_price" bson:"avg_special_price"`
	SumOfStockQuantity int `json:"sum_of_stock_quantity" bson:"sum_of_stock_quantity"`
	MinOfStockQuantity int `json:"min_of_stock_quantity" bson:"min_of_stock_quantity"`
}

func (bamiloCatalogConfigInfo *BamiloCatalogConfigInfo) SetVisibleInShopTrue() {

	bamiloCatalogConfigInfo.VisibleInShopBool = true

}

func (bamiloCatalogConfigInfo *BamiloCatalogConfigInfo) UpsertConfigInfo(session *mgo.Session) {

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("competition_analysis").C("bml_catalog_config")
	_, err := c.UpsertId(bamiloCatalogConfigInfo.IDBmlCatalogConfig, bamiloCatalogConfigInfo)
	if err != nil {
		log.Println("Error creating Profile: ", err.Error())
	}
}

type BamiloCatalogConfigInfoHist struct {
	// qualitative information
	IDBmlCatalogConfigHist int       `json:"id_bml_catalog_config_hist" bson:"id_bml_catalog_config_hist"`
	FKBmlCatalogConfig     int       `json:"fk_bml_catalog_config" bson:"fk_bml_catalog_config"`
	ConfigSnapshotAt       time.Time `json:"config_snapshot_at" bson:"config_snapshot_at"` // for sales: order_created_at

	// historical visibility
	VisibleInShopBool bool `json:"visible_in_shop" bson:"visible_in_shop"`
	// historical price and quantity
	AvgPrice           int `json:"avg_price" bson:"avg_price"`
	AvgSpecialPrice    int `json:"avg_special_price" bson:"avg_special_price"`
	SumOfStockQuantity int `json:"sum_of_stock_quantity" bson:"sum_of_stock_quantity"`
	MinOfStockQuantity int `json:"min_of_stock_quantity" bson:"min_of_stock_quantity"`
}

func (bamiloCatalogConfigInfoHist *BamiloCatalogConfigInfoHist) SetVisibleInShopTrue() {

	bamiloCatalogConfigInfoHist.VisibleInShopBool = true

}

func (bamiloCatalogConfigInfoHist *BamiloCatalogConfigInfoHist) UpsertConfigInfoHist(session *mgo.Session) {

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("competition_analysis").C("bml_catalog_config_hist")
	_, err := c.UpsertId(bamiloCatalogConfigInfoHist.IDBmlCatalogConfigHist, bamiloCatalogConfigInfoHist)
	if err != nil {
		log.Println("Error creating Profile: ", err.Error())
	}

}

// config sales -------------------------------------------------------------------------

type BamiloCatalogConfigSales struct {
	// qualitative information
	IDBmlCatalogConfig int       `json:"id_bml_catalog_config" bson:"id_bml_catalog_config"`
	ConfigSnapshotAt   time.Time `json:"config_snapshot_at" bson:"config_snapshot_at"`
	// sales
	CountOfSoi            int `json:"count_of_soi" bson:"count_of_soi"`
	SumOfUnitPrice        int `json:"sum_of_unit_price" bson:"sum_of_unit_price"`
	SumOfPaidPrice        int `json:"sum_of_paid_price" bson:"sum_of_paid_price"`
	SumOfCouponMoneyValue int `json:"sum_of_coupon_money_value" bson:"sum_of_coupon_money_value"`
	SumOfCartRuleDiscount int `json:"sum_of_cart_rule_discount" bson:"sum_of_cart_rule_discount"`
}

func (bamiloCatalogConfigSales *BamiloCatalogConfigSales) UpsertConfigSales(session *mgo.Session) {

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("competition_analysis").C("bml_catalog_config")
	_, err := c.UpsertId(bamiloCatalogConfigSales.IDBmlCatalogConfig, bamiloCatalogConfigSales)
	if err != nil {
		log.Println("Error creating Profile: ", err.Error())

	}

}

type BamiloCatalogConfigSalesHist struct {
	// qualitative information
	IDBmlCatalogConfigHist int       `json:"id_bml_catalog_config_hist" bson:"id_bml_catalog_config_hist"`
	FKBmlCatalogConfig     int       `json:"fk_bml_catalog_config" bson:"fk_bml_catalog_config"`
	ConfigSnapshotAt       time.Time `json:"config_snapshot_at" bson:"config_snapshot_at"` // for sales: order_created_at
	// sales
	CountOfSoi            int `json:"count_of_soi" bson:"count_of_soi"`
	SumOfUnitPrice        int `json:"sum_of_unit_price" bson:"sum_of_unit_price"`
	SumOfPaidPrice        int `json:"sum_of_paid_price" bson:"sum_of_paid_price"`
	SumOfCouponMoneyValue int `json:"sum_of_coupon_money_value" bson:"sum_of_coupon_money_value"`
	SumOfCartRuleDiscount int `json:"sum_of_cart_rule_discount" bson:"sum_of_cart_rule_discount"`
}

func (bamiloCatalogConfigSalesHist *BamiloCatalogConfigSalesHist) UpsertConfigSalesHist(session *mgo.Session) {

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("competition_analysis").C("bml_catalog_config_hist")
	_, err := c.UpsertId(bamiloCatalogConfigSalesHist.IDBmlCatalogConfigHist, bamiloCatalogConfigSalesHist)
	if err != nil {
		log.Println("Error creating Profile: ", err.Error())

	}

}
