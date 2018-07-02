package bamilocatalogconfig

import (
	"time"
)

type BamiloCatalogConfig struct {
	// qualitative information
	IDBmlCatalogConfig string    `json:"id_bml_catalog_config" bson:"id_bml_catalog_config"`
	FKBmlCatalogConfig string    `json:"fk_bml_catalog_config" bson:"fk_bml_catalog_config"`
	IDMongoDb          string    `json:"_id" bson:"_id"`
	ConfigSnapshotAt   time.Time `json:"config_snapshot_at" bson:"config_snapshot_at"` // for sales: order_created_at
	SKU                string    `json:"sku" bson:"sku"`
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
	VisibleInShop     string
	VisibleInShopBool bool `json:"visible_in_shop" bson:"visible_in_shop"`
	// historical price and quantity
	AvgPrice           string `json:"avg_price" bson:"avg_price"`
	AvgSpecialPrice    string `json:"avg_special_price" bson:"avg_special_price"`
	SumOfStockQuantity string `json:"sum_of_stock_quantity" bson:"sum_of_stock_quantity"`
	MinOfStockQuantity string `json:"min_of_stock_quantity" bson:"min_of_stock_quantity"`
	// sales
	CountOfSoi            string `json:"count_of_soi" bson:"count_of_soi"`
	SumOfUnitPrice        string `json:"sum_of_unit_price" bson:"sum_of_unit_price"`
	SumOfPaidPrice        string `json:"sum_of_paid_price" bson:"sum_of_paid_price"`
	SumOfCouponMoneyValue string `json:"sum_of_coupon_money_value" bson:"sum_of_coupon_money_value"`
	SumOfCartRuleDiscount string `json:"sum_of_cart_rule_discount" bson:"sum_of_cart_rule_discount"`
}

func (bamiloCatalogConfig *BamiloCatalogConfig) SetVisibleInShopTrue() {

	bamiloCatalogConfig.VisibleInShopBool = true

}
