module github.com/videocoin/cloud-streams

go 1.14

require (
	github.com/AlekSi/pointer v1.1.0
	github.com/ethereum/go-ethereum v1.9.11
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.3.3
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway v1.11.3
	github.com/jinzhu/copier v0.0.0-20190924061706-b57f9002281a
	github.com/jinzhu/gorm v1.9.12
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/labstack/gommon v0.3.0 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/opentracing/opentracing-go v1.1.0
	github.com/robertkrimen/otto v0.0.0-20170205013659-6a77b7cbc37d // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/streadway/amqp v0.0.0-20200108173154-1c71cc93ed71
	github.com/videocoin/cloud-api v1.1.6
	github.com/videocoin/cloud-pkg v1.0.0
	google.golang.org/grpc v1.27.1
	gopkg.in/go-playground/validator.v9 v9.31.0
)

replace github.com/videocoin/cloud-api => ../cloud-api

replace github.com/videocoin/cloud-pkg => ../cloud-pkg
