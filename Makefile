NAME := watch

TAG_NAME := $(shell git tag -l --contains HEAD)
SHA := $(shell git rev-parse HEAD)
VERSION := $(if $(TAG_NAME),$(TAG_NAME),$(SHA))
REVISION := $(SHA)
SRCS	:= $(shell find . -type d -name archive -prune -o -type f -name '*.go')
IMAGE ?=
TAG ?=

deps:
	go get -v

build/static: $(SRCS) deps
	NAME=$(NAME) VERSION=$(VERSION) REVISION=$(REVISION) sh ./build.sh

build/docker-alpine:
	docker build -f ci/alpine.dockerfile -t $(IMAGE):$(TAG) .

build/docker-debian:
	docker build -f ci/debian.dockerfile -t $(IMAGE):$(TAG) .

build/docker: build/docker-alpine build/docker-debian