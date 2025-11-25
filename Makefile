.PHONY: start
start:
	CONFIG_PATH=./configs/local.yaml go run ./cmd/url-shortener/main.go

.DEFAULT_GOAL := start