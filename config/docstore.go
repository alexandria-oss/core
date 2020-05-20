package config

import "github.com/spf13/viper"

type docstore struct {
	Collection   string
	PartitionKey string
	SortKey      string
	AllowScan    bool
}

func init() {
	viper.SetDefault("alexandria.persistence.doc.collection", "ALEXANDRIA_COLLECTION")
	viper.SetDefault("alexandria.persistence.doc.partition_key", "")
	viper.SetDefault("alexandria.persistence.doc.sort_key", "")
	viper.SetDefault("alexandria.persistence.doc.allow_scan", false)
}

func newDocstoreConfig() docstore {
	return docstore{
		Collection:   viper.GetString("alexandria.persistence.doc.collection"),
		PartitionKey: viper.GetString("alexandria.persistence.doc.partition_key"),
		SortKey:      viper.GetString("alexandria.persistence.doc.sort_key"),
		AllowScan:    viper.GetBool("alexandria.persistence.doc.allow_scan"),
	}
}
