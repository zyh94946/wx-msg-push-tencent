NAME := main
BUILD_DIR := build
GOBUILD = CGO_ENABLED=0 $(GO_DIR)go build -ldflags="-s -w -buildid=" -o $(BUILD_DIR)

.PHONY: build

build:
	@echo "Building ..."
	mkdir -p $(BUILD_DIR)
	$(GOBUILD)/$(NAME)
	@echo "Build success"

%.zip: %
	@zip -du $(NAME)-$@ -j $(BUILD_DIR)/$</*
	@echo "<<< ---- $(NAME)-$@"

release: linux-amd64.zip

linux-amd64:
	mkdir -p $(BUILD_DIR)/$@
	GOARCH=amd64 GOOS=linux $(GOBUILD)/$@/$(NAME)