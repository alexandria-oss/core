module github.com/alexandria-oss/core

go 1.13

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-kit/kit v0.10.0
	github.com/go-redis/redis/v7 v7.2.0
	github.com/google/uuid v1.1.1
	github.com/gorilla/mux v1.7.3
	github.com/opentracing/opentracing-go v1.1.0
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.4.5
	github.com/openzipkin/zipkin-go v0.2.2
	github.com/prometheus/client_golang v1.3.0
	github.com/rs/cors v1.7.0
	github.com/sony/gobreaker v0.4.1
	github.com/sony/sonyflake v1.0.0
	github.com/spf13/viper v1.6.3
	github.com/stretchr/testify v1.5.1
	go.uber.org/zap v1.13.0
	gocloud.dev v0.19.0
	gocloud.dev/pubsub/kafkapubsub v0.19.0
	golang.org/x/crypto v0.0.0-20200206161412-a0c6ece9d31a
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	google.golang.org/grpc v1.27.1
)
