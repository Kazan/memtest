.PHONY: build

export GO111MODULE=on
export CGO_ENABLED=0


client:
	go run client/main.go

server:
	go run server/main.go
