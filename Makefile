HOSTNAME?=github.com
NAMESPACE?=zalopay-oss
NAME?=nacos
BINARY=terraform-provider-${NAME}
VERSION?=0.1.0
OS_ARCH?=linux_amd64

build:
	go build -o bin/${BINARY}

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv bin/${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

testacc:
	TF_ACC=1 go test -v -coverprofile=coverage.out ./internal/...
