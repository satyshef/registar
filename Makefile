APP_PATH="registar"

.PHONY: build
build:
	go build -modfile go.mod -v -o ${APP_PATH} ./cmd/app

.DEFAULT_GOAL := build
