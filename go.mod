module github.com/sniperCore

go 1.16

require (
	github.com/Shopify/sarama v1.36.0
	github.com/aws/aws-sdk-go v1.43.21
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.8.1
	github.com/go-redis/redis/v8 v8.11.5
	github.com/hashicorp/consul/api v1.15.2
	github.com/jonboulle/clockwork v0.3.0 // indirect
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/olivere/elastic/v7 v7.0.32
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
	github.com/sirupsen/logrus v1.9.0
	github.com/spf13/viper v1.13.0
	google.golang.org/grpc v1.46.2
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.2.0 // indirect
	google.golang.org/protobuf v1.28.1
	gopkg.in/gographics/imagick.v3 v3.4.2
	gorm.io/driver/clickhouse v0.2.1
	gorm.io/driver/mysql v1.1.1
	gorm.io/driver/postgres v1.1.0
	gorm.io/driver/sqlite v1.1.4
	gorm.io/driver/sqlserver v1.0.7
	gorm.io/gorm v1.21.11
)
