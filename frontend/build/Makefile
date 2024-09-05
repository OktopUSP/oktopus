.PHONY: help build push start stop release remove delete run logs bash

DOCKER_USER ?= oktopusp
DOCKER_APP ?= frontend
DOCKER_TAG ?= $(shell git log --format="%h" -n 1)
CONTAINER_SHELL ?= /bin/sh

.DEFAULT_GOAL := help

help:
	@echo "Makefile arguments:"
	@echo ""
	@echo "DOCKER_USER - docker user to build image"
	@echo "DOCKER_APP - docker image name"
	@echo "DOCKER_TAG - docker image tag"
	@echo "CONTAINER_SHELL - container shell e.g:'/bin/bash'"
	@echo ""
	@echo "Makefile commands:"
	@echo ""
	@echo "build - docker image build"
	@echo "release - tag image as latest and push to registry"

build:
	@docker build -t ${DOCKER_USER}/${DOCKER_APP}:${DOCKER_TAG} -f Dockerfile ../

release: build
	@docker push ${DOCKER_USER}/${DOCKER_APP}:${DOCKER_TAG}
	@docker tag ${DOCKER_USER}/${DOCKER_APP}:${DOCKER_TAG} ${DOCKER_USER}/${DOCKER_APP}:latest
	@docker push ${DOCKER_USER}/${DOCKER_APP}:latest