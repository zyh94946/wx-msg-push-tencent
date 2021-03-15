.PHONY: build
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

build:
	@echo "Building ..."
	@go build -o main
	@upx main
	@echo "Build success"
