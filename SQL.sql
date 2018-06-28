	SELECT 

	cbchv.id_bml_catalog_config_hist

	FROM baa_application.commercial.cpt_bml_config_hist_view cbchv
	WHERE cbchv.id_bml_catalog_config_hist = CONCAT(171499, REPLACE(CONVERT (CHAR(10), GETDATE(), 101),'/',''));

  DELETE FROM baa_application.commercial.cpt_bml_catalog_config