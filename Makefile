.PHONY: build arm64 clean npm-install gen-webui backend
GO=GOPROXY="https://goproxy.io" CGO_ENABLED=1 GO111MODULE=on go
GOARM64=GOPROXY="https://goproxy.io" CGO_ENABLED=1 GO111MODULE=on GOOS=linux GOARCH=arm64  go
BIN=./edge_admin
clean:
	rm -rf ./edge_admin

npm-install:
	cd webui && npm install && cd ..

gen-webui:
	cd webui && npm run build:stage && cd ..
	
build: clean gen-webui
	go generate
	$(GO) build -o $(BIN) .

arm64:
	$(GOARM64) build -o $(BIN) . 

backend:
	$(GO) build -o $(BIN) .
