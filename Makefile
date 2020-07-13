GOOS?=linux
GOARCH?=amd64

GCP_PROJECT?=videocoin-network

NAME=streams
VERSION=$$(git describe --abbrev=0)-$$(git rev-parse --abbrev-ref HEAD)-$$(git rev-parse --short HEAD)

DBM_MSQLURI=root:@tcp(127.0.0.1:3306)/videocoin?charset=utf8&parseTime=True&loc=Local
ENV?=dev

.PHONY: deploy

default: build

version:
	@echo ${VERSION}

build:
	GOOS=${GOOS} GOARCH=${GOARCH} \
		go build -mod vendor \
			-ldflags="-w -s -X main.Version=${VERSION}" \
			-o bin/${NAME} \
			./cmd/main.go

deps:
	GO111MODULE=on go mod vendor

lint: docker-lint

docker-lint:
	docker build -f Dockerfile.lint .

docker-build:
	docker build -t gcr.io/${GCP_PROJECT}/${NAME}:${VERSION} -f Dockerfile .

docker-push:
	docker push gcr.io/${GCP_PROJECT}/${NAME}:${VERSION}

dbm-status:
	./tools/${GOOS}/goose -dir migrations -table ${NAME} mysql "${DBM_MSQLURI}" status

dbm-up:
	./tools/${GOOS}/goose -dir migrations -table ${NAME} mysql "${DBM_MSQLURI}" up

dbm-down:
	./tools/${GOOS}/goose -dir migrations -table ${NAME} mysql "${DBM_MSQLURI}" down

release: docker-build docker-push

deploy:
	ENV=${ENV} GCP_PROJECT=${GCP_PROJECT} deploy/deploy.sh
