.PHONY: all tidy build test run

all: tidy build test

TAGS?=

MOD_DIR=control-plane

 tidy:
	cd $(MOD_DIR) && go mod tidy

 build:
	cd $(MOD_DIR) && go build ./...

 test:
	cd $(MOD_DIR) && go test ./...

 run:
	cd $(MOD_DIR) && go run ./cmd/api