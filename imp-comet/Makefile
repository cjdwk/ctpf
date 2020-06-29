VERSION := $(shell echo $(shell cat version | grep "Version" | cut -d '=' -f2))
APP_NAME := ctpf-comet
BUILD_DIR := "build"
WORK_DIR := cmd
APP := ${BUILD_DIR}/${APP_NAME}
TARGET := 172.16.103.31:/root/${APP_NAME}/bin
PASSWORD := "123456"

build:
	GOOS=linux GOARCH=amd64 GO111MODULE=on go build -v -o $(APP).$(VERSION) $(WORK_DIR)/*.go
	cp etc/config.toml ${BUILD_DIR}

run:

upload:
	sshpass -p $(PASSWORD) scp $(APP).$(VERSION) root@$(TARGET)

ver:
	@echo $(VERSION)

docker-run:
	docker build -t chat33-service .
	docker run -p 50051:50051 chat33-service

clean: ## Remove previous build
	@rm -rf ${BUILD_DIR}
	@go clean