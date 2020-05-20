package config

import "github.com/spf13/viper"

type aws struct {
	CognitoPoolID   string
	CognitoClientID string
}

func init() {
	viper.SetDefault("alexandria.cloud.aws.cognito.pool", "example_pool_id")
	viper.SetDefault("alexandria.cloud.aws.cognito.client", "example_client_id")
}

func newAWSConfig() aws {
	return aws{
		CognitoPoolID:   viper.GetString("alexandria.cloud.aws.cognito.pool"),
		CognitoClientID: viper.GetString("alexandria.cloud.aws.cognito.client"),
	}
}
