
# https://stackoverflow.com/questions/34712972/in-a-makefile-how-can-i-fetch-and-assign-a-git-commit-hash-to-a-variable
GIT_COMMIT_HASH=$(shell git rev-parse --short HEAD)
SWAG_VERSION ?= v2.0.0-rc5

user-service-build:
	docker build \
		--build-arg GIT_COMMIT_HASH=${GIT_COMMIT_HASH} \
		--tag totocorpsoftwareinc/user-service:${GIT_COMMIT_HASH} \
		-f build/user-service/Dockerfile \
		.

generate-api-spec:
	cd cmd/users && \
	go run github.com/swaggo/swag/v2/cmd/swag@${SWAG_VERSION} init \
		--v3.1 \
		--generalInfo main.go \
		--dir .,../../internal/controller,../../pkg/communication \
		--output ../../api \
		--outputTypes go,yaml \
		--parseDependency \
		--parseInternal \
		--generatedTime=false
