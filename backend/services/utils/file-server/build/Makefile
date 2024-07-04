.PHONY: help build push start stop release remove delete run logs bash

DOCKER_USER ?= oktopusp
DOCKER_APP ?= file-server
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
	@echo "push - push docker iamge to registry"
	@echo "run - create and start docker container with the image"
	@echo "start - start existent docker container with the image"
	@echo "stop - stop docker container running the image"
	@echo "remove - remove docker container running the image"
	@echo "delete - delete docker image"
	@echo "logs - show logs of docker container"
	@echo "bash - access container shell"
	@echo "release - tag image as latest and push to registry"

build:
	@docker build -t ${DOCKER_USER}/${DOCKER_APP}:${DOCKER_TAG} -f Dockerfile ../

run:
	@docker run -d --name ${DOCKER_USER}-${DOCKER_APP} ${DOCKER_USER}/${DOCKER_APP}:${DOCKER_TAG}

stop:
	@docker stop ${DOCKER_USER}-${DOCKER_APP}

remove: stop
	@docker rm ${DOCKER_USER}-${DOCKER_APP}

delete:
	@docker rmi ${DOCKER_USER}/${DOCKER_APP}:${DOCKER_TAG}

start:
	@docker start ${DOCKER_USER}-${DOCKER_APP}

push:
	@docker push ${DOCKER_USER}/${DOCKER_APP}:${DOCKER_TAG}

logs:
	@docker logs -f ${DOCKER_USER}-${DOCKER_APP}

bash:
	@docker exec -it ${DOCKER_USER}-${DOCKER_APP} ${CONTAINER_SHELL}

release: build
	@docker push ${DOCKER_USER}/${DOCKER_APP}:${DOCKER_TAG}
	@docker tag ${DOCKER_USER}/${DOCKER_APP}:${DOCKER_TAG} ${DOCKER_USER}/${DOCKER_APP}:latest
	@docker push ${DOCKER_USER}/${DOCKER_APP}:latest