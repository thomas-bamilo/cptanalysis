SELECT COUNT(cdchv.fk_dgk_catalog_config) FROM commercial.cpt_dgk_config_hist_view cdchv


-- DIGIKALA TABLES

CREATE TABLE baa_application.commercial.cpt_dgk_catalog_config (
  id_dgk_catalog_config INT IDENTITY(1,1) PRIMARY KEY
  ,config_updated_at DATE NOT NULL DEFAULT (GETDATE())
  ,sku_name NVARCHAR(200) NOT NULL UNIQUE
  ,img_link VARCHAR(300)
  ,sku_rank INT
);

--Create a rowstore table with a unique constraint.  
--The unique constraint is implemented as a nonclustered index.  
CREATE TABLE baa_application.commercial.cpt_dgk_catalog_config_hist (
  id_dgk_catalog_config_hist AS CONCAT(fk_dgk_catalog_config,REPLACE(CONVERT(VARCHAR(50), config_snapshot_at, 101),'/',''))
  ,fk_dgk_catalog_config INT NOT NULL FOREIGN KEY REFERENCES baa_application.commercial.cpt_dgk_catalog_config(id_dgk_catalog_config)
  ,config_snapshot_at DATE NOT NULL DEFAULT (GETDATE())
  ,rating FLOAT
  ,avg_price FLOAT
  ,avg_special_price FLOAT
  ,is_out_of_stock SMALLINT
  CONSTRAINT dgk_catalog_config_hist_check UNIQUE (fk_dgk_catalog_config,config_snapshot_at)  
);  

--Store the table as a columnstore.   
--The unique constraint is preserved as a nonclustered index on the columnstore table.  
CREATE CLUSTERED COLUMNSTORE INDEX dgk_catalog_config_hist_cci ON baa_application.commercial.cpt_dgk_catalog_config_hist;

  -- if you want to drop a column
-- ALTER TABLE baa_application.commercial.cpt_dgk_catalog_config_hist DROP COLUMN id_dgk_catalog_config_hist

--By using the previous two steps, every row in the table meets the UNIQUE constraint  
--on a non-NULL column.  
--This has the same end-result as having a primary key constraint  
--All updates and inserts must meet the unique constraint on the nonclustered index or they will fail.  

--If desired, add a foreign key constraint on fk_dgk_catalog_config: no need: it was already declared in the intitial table!

-- ALTER TABLE baa_application.commercial.cpt_dgk_catalog_config_hist  
-- WITH CHECK ADD FOREIGN KEY(fk_dgk_catalog_config) REFERENCES baa_application.commercial.cpt_dgk_catalog_config(id_dgk_catalog_config);

DELETE FROM baa_application.commercial.cpt_dgk_catalog_config_hist
DELETE FROM baa_application.commercial.cpt_dgk_catalog_config



  -- BAMILO TABLES

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


CREATE TABLE baa_application.commercial.cpt_bml_catalog_config_hist (
  id_bml_catalog_config_hist AS CONCAT(fk_bml_catalog_config,REPLACE(CONVERT(VARCHAR(50), config_snapshot_at, 101),'/',''))
  ,fk_bml_catalog_config INT NOT NULL FOREIGN KEY REFERENCES baa_application.commercial.cpt_bml_catalog_config(id_bml_catalog_config)
  ,config_snapshot_at DATE NOT NULL DEFAULT (GETDATE())

  ,visible_in_shop SMALLINT

  ,avg_price FLOAT
  ,avg_special_price FLOAT
  ,sum_of_stock_quantity INT
  ,min_of_stock_quantity INT

  ,count_of_soi INT
  ,sum_of_unit_price FLOAT
  ,sum_of_paid_price FLOAT
  ,sum_of_coupon_money_value FLOAT
  ,sum_of_cart_rule_discount FLOAT

  CONSTRAINT bml_catalog_config_hist_check UNIQUE (fk_bml_catalog_config,config_snapshot_at)

);

CREATE CLUSTERED COLUMNSTORE INDEX bml_catalog_config_hist_cci ON baa_application.commercial.cpt_bml_catalog_config_hist;





CREATE TABLE baa_application.commercial.cpt_bml_dgk_catalog_config (

  fk_bml_catalog_config INT NOT NULL FOREIGN KEY REFERENCES baa_application.commercial.cpt_bml_catalog_config(id_bml_catalog_config)
  ,fk_dgk_catalog_config INT NOT NULL FOREIGN KEY REFERENCES baa_application.commercial.cpt_dgk_catalog_config(id_dgk_catalog_config)

  CONSTRAINT id_bml_dgk_catalog_config PRIMARY KEY (fk_bml_catalog_config,fk_dgk_catalog_config)

);

UPDATE baa_application.commercial.cpt_bml_catalog_config
SET 
			fk_bml_catalog_config

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


  INSERT INTO  baa_application.commercial.cpt_dgk_catalog_config_hist (
    fk_dgk_catalog_config
    ,rating
    ,avg_price
    ,avg_special_price
    ,is_out_of_stock)
  VALUES (@id_dgk_catalog_config, 
          @rating, 
          @avg_price, 
          @avg_special_price,
          @is_out_of_stock);
