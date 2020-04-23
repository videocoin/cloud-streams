module github.com/videocoin/cloud-streams

go 1.13

replace github.com/videocoin/cloud-pkg => ../cloud-pkg

replace github.com/videocoin/cloud-api => ../cloud-api

require (
	github.com/AlekSi/pointer v1.1.0
	github.com/ethereum/go-ethereum v1.9.10
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.3.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0
	github.com/jinzhu/copier v0.0.0-20190924061706-b57f9002281a
	github.com/jinzhu/gorm v1.9.12
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/labstack/gommon v0.3.0 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/opentracing/opentracing-go v1.1.0
	github.com/sirupsen/logrus v1.4.2
	github.com/streadway/amqp v0.0.0-20200108173154-1c71cc93ed71
	github.com/videocoin/cloud-api v0.0.17
	github.com/videocoin/cloud-pkg v0.0.5
	google.golang.org/grpc v1.26.0
	gopkg.in/go-playground/validator.v9 v9.31.0
)
