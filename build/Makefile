GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOLIST=$(GOCMD) list
GODEP=godep save
VERSION?=

APP_NAME?=
WORKSPACE?=$(GOPATH)/src/dinit
BIN_PATH?=$(WORKSPACE)/build/bin
SOURCE_PATH=$(WORKSPACE)/pkg
GIT_COMMIT=$(shell git rev-parse --short HEAD)
BUILD_TIME=$(shell date "+%Y-%m-%dT%H:%M:%S%z")
MD5SUM=$(shell md5 -q $(BIN_PATH)/$(APP_NAME))
LDFLAGS += -X $(shell $(GOLIST) $(SOURCE_PATH)/$(APP_NAME)/app).GitHash=$(GIT_COMMIT)
LDFLAGS += -X $(shell $(GOLIST) $(SOURCE_PATH)/$(APP_NAME)/app).Version=$(VERSION)
LDFLAGS += -X $(shell $(GOLIST) $(SOURCE_PATH)/$(APP_NAME)/app).MD5SUM=$(MD5SUM)
GOFILES = $(shell find $(WORKSPACE)/pkg/$(APP_NAME) -type f -name '*.go')
GOFILES_ALL = $(GOFILES)

${APP_NAME}: $(GOFILES_ALL)
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -v -ldflags "$(LDFLAGS)" -o $(BIN_PATH)/$(APP_NAME)
