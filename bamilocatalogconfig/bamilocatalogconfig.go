package bamilocatalogconfig

type BamiloCatalogConfig struct {
	// qualitative information
	IDBmlCatalogConfig string `json:"id_bml_catalog_config"`
	ConfigSnapshotAt   string `json:"config_snapshot_at"` // for sales: order_created_at
	SKU                string `json:"sku"`
	SKUName            string `json:"sku_name"`
	ImgLink            string `json:"img_link"`
	Description        string `json:"description"`
	ShortDescription   string `json:"short_description"`
	PackageContent     string `json:"package_content"`
	ProductWarranty    string `json:"product_warranty"`
	// category
	FKBiCategoryOne     string `json:"fk_bi_category_one"`
	BiCategoryOneName   string `json:"bi_category_one_name"`
	FKBiCategoryTwo     string `json:"fk_bi_category_two"`
	BiCategoryTwoName   string `json:"bi_category_two_name"`
	FKBiCategoryThree   string `json:"fk_bi_category_three"`
	BiCategoryThreeName string `json:"bi_category_three_name"`
	// brand
	FKCatalogBrand string `json:"fk_catalog_brand"`
	BrandName      string `json:"brand_name"`
	BrandNameEn    string `json:"brand_name_en"`
	BrandStatus    string `json:"brand_status"`
	// supplier
	FkSupplier     string `json:"fk_supplier"`
	SupplierName   string `json:"supplier_name"`
	SupplierNameEn string `json:"supplier_name_en"`
	SupplierStatus string `json:"supplier_status"`
	// historical visibility
	VisibleInShop string `json:"visible_in_shop"`
	// historical price and quantity
	AvgPrice           string `json:"avg_price"`
	AvgSpecialPrice    string `json:"avg_special_price"`
	SumOfStockQuantity string `json:"sum_of_stock_quantity"`
	MinOfStockQuantity string `json:"min_of_stock_quantity"`
	// sales
	CountOfSoi            string `json:"count_of_soi"`
	SumOfUnitPrice        string `json:"sum_of_unit_price"`
	SumOfPaidPrice        string `json:"sum_of_paid_price"`
	SumOfCouponMoneyValue string `json:"sum_of_coupon_money_value"`
	SumOfCartRuleDiscount string `json:"sum_of_cart_rule_discount"`
}
