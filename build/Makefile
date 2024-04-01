.PHONY: help build

.DEFAULT_GOAL := help
DOCKER_USER ?= oktopus

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
	@echo "build - docker images build backend+frontend"
	@echo "build-backend - docker images build backend"
	@echo "build-frontend - docker images build frontend"
	@echo "release - new docker images latest and push to docker hub"
	@echo "release-backend - new backend docker images latest and push to docker hub"
	@echo "release-frontend - new frontend docker images latest and push to docker hub"

build: build-frontend build-backend

build-backend:
	@make build -C ../backend/services/controller/build/ DOCKER_USER=${DOCKER_USER}
	@make build -C ../backend/services/utils/socketio/build/ DOCKER_USER=${DOCKER_USER}
	@make build -C ../backend/services/mtp/adapter/build/ DOCKER_USER=${DOCKER_USER}
	@make build -C ../backend/services/mtp/ws-adapter/build/ DOCKER_USER=${DOCKER_USER}
	@make build -C ../backend/services/mtp/ws/build/ DOCKER_USER=${DOCKER_USER}
	@make build -C ../backend/services/mtp/mqtt-adapter/build/ DOCKER_USER=${DOCKER_USER}
	@make build -C ../backend/services/mtp/mqtt/build/ DOCKER_USER=${DOCKER_USER}
	@make build -C ../backend/services/mtp/stomp-adapter/build/ DOCKER_USER=${DOCKER_USER}
	@make build -C ../backend/services/mtp/stomp/build/ DOCKER_USER=${DOCKER_USER}

build-frontend:
	@make build -C ../frontend/build/ DOCKER_USER=${DOCKER_USER}

release: release-frontend release-backend

release-backend:
	@make release -C ../backend/services/controller/build/ DOCKER_USER=${DOCKER_USER}
	@make release -C ../backend/services/utils/socketio/build/ DOCKER_USER=${DOCKER_USER}
	@make release -C ../backend/services/mtp/adapter/build/ DOCKER_USER=${DOCKER_USER}
	@make release -C ../backend/services/mtp/ws-adapter/build/ DOCKER_USER=${DOCKER_USER}
	@make release -C ../backend/services/mtp/ws/build/ DOCKER_USER=${DOCKER_USER}
	@make release -C ../backend/services/mtp/mqtt-adapter/build/ DOCKER_USER=${DOCKER_USER}
	@make release -C ../backend/services/mtp/mqtt/build/ DOCKER_USER=${DOCKER_USER}
	@make release -C ../backend/services/mtp/stomp-adapter/build/ DOCKER_USER=${DOCKER_USER}
	@make release -C ../backend/services/mtp/stomp/build/ DOCKER_USER=${DOCKER_USER}

release-frontend:
	@make release -C ../frontend/build/ DOCKER_USER=${DOCKER_USER}