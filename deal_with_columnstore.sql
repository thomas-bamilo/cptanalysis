-- last snapshot view
CREATE VIEW commercial.cpt_dgk_config_last_snapshot_view AS
SELECT 
  
  cdcc.config_updated_at
  ,cdcc.sku_name
  ,cdcc.img_link
  ,cdcc.sku_rank
  ,last_snapshot.rating
  ,last_snapshot.avg_price
  ,last_snapshot.avg_special_price 
  ,last_snapshot.is_out_of_stock
  
  
FROM

  commercial.cpt_dgk_catalog_config cdcc

JOIN

(SELECT 

cdcch.fk_dgk_catalog_config
,MAX(cdcch.rating) rating
,MAX(cdcch.avg_price) avg_price
,MAX(cdcch.avg_special_price) avg_special_price
,MAX(cdcch.is_out_of_stock) is_out_of_stock
,ROW_NUMBER() OVER(PARTITION BY cdcch.fk_dgk_catalog_config
					   ORDER BY cdcch.config_snapshot_at  DESC) AS rk
  
FROM commercial.cpt_dgk_catalog_config_hist cdcch

GROUP BY cdcch.fk_dgk_catalog_config, cdcch.config_snapshot_at) last_snapshot
ON last_snapshot.fk_dgk_catalog_config = cdcc.id_dgk_catalog_config
AND last_snapshot.rk = 1;


 
-- history view
CREATE VIEW commercial.cpt_dgk_config_hist_view AS
SELECT 
  cdcch.fk_dgk_catalog_config
  ,cdcc.sku_name
  ,cdcc.img_link
  ,cdcc.sku_rank
  ,cdcch.config_snapshot_at
  ,cdcch.rating
  ,cdcch.avg_price
  ,cdcch.avg_special_price
  ,cdcch.is_out_of_stock
  
  
FROM commercial.cpt_dgk_catalog_config cdcc

JOIN commercial.cpt_dgk_catalog_config_hist cdcch

ON  cdcch.fk_dgk_catalog_config = cdcc.id_dgk_catalog_config;


-- select only certain records from the history view
SELECT 

  cdchv.sku_name
  ,cdchv.config_snapshot_at
  ,cdchv.avg_price


FROM commercial.cpt_dgk_config_hist_view cdchv

WHERE cdchv.config_snapshot_at > DATEADD(DAY,-3,GETDATE())
AND cdchv.id_dgk_catalog_config = 3



      