package config

import "github.com/spf13/viper"

type aws struct {
	CognitoPoolID string
}

func init() {
	viper.SetDefault("alexandria.cloud.aws.cognito.pool_id", "example_pool_id")
}

func newAWSConfig() aws {
	return aws{
		CognitoPoolID: viper.GetString("alexandria.cloud.aws.cognito.pool_id"),
	}
}
