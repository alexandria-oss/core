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
func NewDynamoDBCollectionPool(ctx context.Context, cfg *config.Kernel) (*docstore.Collection, func(), error) {
	URL := fmt.Sprintf("dynamodb://%s?partition_key=%s", cfg.Docstore.Collection,
		strings.ToLower(cfg.Docstore.PartitionKey))

	if cfg.Docstore.AllowScan {
		URL += "&allow_scans=true"
	}
	if cfg.Docstore.SortKey != "" {
		URL += "&sort_key=" + strings.ToLower(cfg.Docstore.SortKey)
	}

	db, err := docstore.OpenCollection(ctx, URL)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		_ = db.Close()
	}

	return db, cleanup, nil
}
