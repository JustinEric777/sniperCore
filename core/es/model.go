package es

import (
	"context"

	"github.com/olivere/elastic/v7"
)

var ctx context.Context

func CreateIndex(client *elastic.Client, index string, mapping string) error {
	exists, err := client.IndexExists(index).Do(ctx)
	if err != nil {
		return err
	}

	if exists {
		createIndex, err := client.CreateIndex(index).Body(mapping).Do(ctx)
		if err != nil {
			return err
		}

		if !createIndex.Acknowledged {
			return err
		}
	}

	return nil
}

/**
 * delete Index
 */
func DeleteIndex(client *elastic.Client, index string) error {
	deleteIndex, err := client.DeleteIndex(index).Do(ctx)
	if err != nil {
		return err
	}

	if !deleteIndex.Acknowledged {
		return err
	}

	return nil
}

/**
 * add index data
 */
func AddIndexData(client *elastic.Client, index string, jsonData string) error {
	_, err := client.Index().Index(index).BodyString(jsonData).Do(ctx)
	if err != nil {
		return err
	}

	return nil
}
