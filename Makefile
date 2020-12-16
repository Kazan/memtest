.PHONY: server client req

export GO111MODULE=on
export CGO_ENABLED=0


client:
	go run client/main.go

server:
	go run server/main.go

req:
	curl -v -0 -XGET http://127.0.0.1:4040/
