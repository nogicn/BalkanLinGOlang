.PHONY: run deploy build

run:
	air

build:
	go build -o BalkanLinGO
	CC=aarch64-linux-gnu-gcc CGO_ENABLED=1 GOARCH=arm64 GOOS=linux go build -o BalkanLinGO_arm64 
