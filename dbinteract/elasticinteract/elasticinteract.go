package elasticinteract

import (
	"context"
	"log"

	"gopkg.in/olivere/elastic.v5"
)

func AddBamiloCatalogConfigToIndex(elasticClient *elastic.Client, ctx context.Context, bamiloCatalogConfigElasticJSON string) {

	_, err := elasticClient.Index().
		Index(`bamilo_vs_competitor_sku_history`).
		Type(`bamilo_vs_competitor_sku_history`).
		BodyJson(bamiloCatalogConfigElasticJSON).
		Do(ctx)

	checkError(err)

}

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
