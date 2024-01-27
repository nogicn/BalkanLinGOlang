.PHONY: run deploy build

run:
	air

build:
	go build -o build/BalkanLinGO
	CC=aarch64-linux-gnu-gcc CGO_ENABLED=1 GOARCH=arm64 GOOS=linux go build -o build/BalkanLinGO_arm64
