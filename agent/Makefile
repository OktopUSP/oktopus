.PHONY: help build release

DOCKER_USER ?= oktopusp
DOCKER_APP ?= obuspa
DOCKER_TAG ?= $(shell git log --format="%h" -n 1)
GIT_REPO ?= https://github.com/BroadbandForum/obuspa

.DEFAULT_GOAL := help

help:
	@echo "Makefile arguments:"
	@echo ""
	@echo "DOCKER_USER - docker user to build image"
	@echo "DOCKER_APP - docker image name"
	@echo "DOCKER_TAG - docker image tag"
	@echo ""
	@echo "Makefile commands:"
	@echo ""
	@echo "build - docker image build"
	@echo "release - tag image as latest and push to registry"

build:
	if [ -d "obuspa" ]; then \
		git -C obuspa pull; \
	else \
		git clone ${GIT_REPO} ; \
	fi
	@docker build -t ${DOCKER_USER}/${DOCKER_APP}:${DOCKER_TAG} -f Dockerfile .

release: build 
	@docker push ${DOCKER_USER}/${DOCKER_APP}:${DOCKER_TAG}
	@docker tag ${DOCKER_USER}/${DOCKER_APP}:${DOCKER_TAG} ${DOCKER_USER}/${DOCKER_APP}:latest
	@docker push ${DOCKER_USER}/${DOCKER_APP}:latest