package config

import "github.com/spf13/viper"

type auth struct {
	JWTSecret string
}

func init() {
	viper.SetDefault("alexandria.security.auth.jwt.secret_key", "EXAMPLE_SECRET")
}

func newAuthConfig() auth {
	return auth{
		JWTSecret: viper.GetString("alexandria.security.auth.jwt.secret_key"),
	}
}
