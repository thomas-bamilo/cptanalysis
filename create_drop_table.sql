-- DIGIKALA TABLES

SELECT COUNT(cdchv.fk_dgk_catalog_config) FROM commercial.cpt_dgk_config_hist_view cdchv

CREATE TABLE baa_application.commercial.cpt_dgk_catalog_config (
  id_dgk_catalog_config INT IDENTITY(1,1) PRIMARY KEY
  ,config_updated_at DATE NOT NULL DEFAULT (GETDATE())
  ,sku_name NVARCHAR(200) NOT NULL UNIQUE
  ,img_link VARCHAR(300)
  ,sku_rank INT
);

CREATE NONCLUSTERED INDEX dgk_catalog_config_name_nci ON baa_application.commercial.cpt_dgk_catalog_config (sku_name); -- to check faster if config already exists

CREATE NONCLUSTERED INDEX dgk_catalog_config_date_nci ON baa_application.commercial.cpt_dgk_catalog_config (config_updated_at); -- to check faster if config was updated today

--Create a rowstore table with a unique constraint.  
--The unique constraint is implemented as a nonclustered index.  
CREATE TABLE baa_application.commercial.cpt_dgk_catalog_config_hist (
  id_dgk_catalog_config_hist BIGINT NOT NULL UNIQUE -- to INPUT
  ,fk_dgk_catalog_config INT NOT NULL FOREIGN KEY REFERENCES baa_application.commercial.cpt_dgk_catalog_config(id_dgk_catalog_config)
  ,config_snapshot_at DATE NOT NULL DEFAULT (GETDATE())
  ,rating FLOAT
  ,avg_price FLOAT
  ,avg_special_price FLOAT
  ,is_out_of_stock SMALLINT
  CONSTRAINT dgk_catalog_config_hist_check UNIQUE (fk_dgk_catalog_config,config_snapshot_at)  
);  

CREATE CLUSTERED COLUMNSTORE INDEX dgk_catalog_config_hist_cci ON baa_application.commercial.cpt_dgk_catalog_config_hist;

CREATE UNIQUE INDEX dgk_catalog_config_hist_unique_nci ON baa_application.commercial.cpt_dgk_catalog_config_hist(id_dgk_catalog_config_hist);



  -- BAMILO TABLES

SELECT COUNT(cdchv.fk_bml_catalog_config) FROM commercial.cpt_bml_config_hist_view cdchv

CREATE TABLE baa_application.commercial.cpt_bml_catalog_config (
  id_bml_catalog_config INT PRIMARY KEY
  ,config_updated_at DATE NOT NULL DEFAULT (GETDATE())
  ,sku VARCHAR(100) NOT NULL UNIQUE
  ,sku_name NVARCHAR(200)
  ,img_link VARCHAR(300)
  ,description NVARCHAR(MAX)
  ,short_description NVARCHAR(MAX)
  ,package_content NVARCHAR(MAX)
  ,product_warranty NVARCHAR(MAX)

  ,bi_category_one_name VARCHAR(50)
  ,bi_category_two_name VARCHAR(50)
  ,bi_category_three_name VARCHAR(50)

  ,brand_name NVARCHAR(100)
  ,brand_name_en NVARCHAR(100)

  ,supplier_name NVARCHAR(100)
  ,supplier_name_en NVARCHAR(100)
);

CREATE NONCLUSTERED INDEX bml_catalog_config_date_nci ON baa_application.commercial.cpt_bml_catalog_config(config_updated_at); -- to check faster if config was updated today


CREATE TABLE baa_application.commercial.cpt_bml_catalog_config_hist (
  id_bml_catalog_config_hist BIGINT NOT NULL UNIQUE -- to INPUT
  ,fk_bml_catalog_config INT NOT NULL FOREIGN KEY REFERENCES baa_application.commercial.cpt_bml_catalog_config(id_bml_catalog_config)
  ,config_snapshot_at DATE NOT NULL DEFAULT (GETDATE())

  ,visible_in_shop SMALLINT

  ,avg_price FLOAT 
  ,avg_special_price FLOAT
  ,sum_of_stock_quantity BIGINT
  ,min_of_stock_quantity BIGINT

  ,count_of_soi BIGINT DEFAULT(0)
  ,sum_of_unit_price FLOAT DEFAULT(0)
  ,sum_of_paid_price FLOAT DEFAULT(0)
  ,sum_of_coupon_money_value FLOAT DEFAULT(0)
  ,sum_of_cart_rule_discount FLOAT DEFAULT(0)

  CONSTRAINT bml_catalog_config_hist_check UNIQUE (fk_bml_catalog_config,config_snapshot_at)

);

CREATE CLUSTERED COLUMNSTORE INDEX bml_catalog_config_hist_cci ON baa_application.commercial.cpt_bml_catalog_config_hist;

CREATE UNIQUE INDEX bml_catalog_config_hist_unique_nci ON baa_application.commercial.cpt_bml_catalog_config_hist(id_bml_catalog_config_hist); -- only search and update on this!





-- Digikala AND Bamilo
CREATE TABLE baa_application.commercial.cpt_bml_dgk_catalog_config (

  fk_bml_catalog_config INT NOT NULL FOREIGN KEY REFERENCES baa_application.commercial.cpt_bml_catalog_config(id_bml_catalog_config)
  ,fk_dgk_catalog_config INT NOT NULL FOREIGN KEY REFERENCES baa_application.commercial.cpt_dgk_catalog_config(id_dgk_catalog_config)

  CONSTRAINT id_bml_dgk_catalog_config PRIMARY KEY (fk_bml_catalog_config,fk_dgk_catalog_config)

);