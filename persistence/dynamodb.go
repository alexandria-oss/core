package persistence

import (
	"context"
	"fmt"
	"strings"

	"github.com/alexandria-oss/core/config"
	"gocloud.dev/docstore"
	_ "gocloud.dev/docstore/awsdynamodb"
)

// NewDynamoDBCollectionPool Obtain an AWS DynamoDB collection connection pool
func NewDynamoDBCollectionPool(ctx context.Context, cfg *config.KernelConfiguration) (*docstore.Collection, func(), error) {
	URL := fmt.Sprintf("dynamodb://%s?partition_key=%s", cfg.DocstoreConfig.Collection,
		strings.ToLower(cfg.DocstoreConfig.PartitionKey))

	if cfg.DocstoreConfig.AllowScan {
		URL += "&allow_scans=true"
	}
	if cfg.DocstoreConfig.SortKey != "" {
		URL += "&sort_key=" + strings.ToLower(cfg.DocstoreConfig.SortKey)
	}

	db, err := docstore.OpenCollection(ctx, URL)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		err = db.Close()
	}

	return db, cleanup, nil
}
