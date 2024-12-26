DOMAIN ?= 
KEY ?=
CLIENT_ID ?=
REDIRECT_URI ?= http://localhost:8080/auth/callback
KEY_FILE ?=

deps-update:
	go mod tidy

run:
	DOMAIN=$(DOMAIN) KEY=$(KEY) CLIENT_ID=$(CLIENT_ID) REDIRECT_URI=$(REDIRECT_URI) KEY_FILE=$(KEY_FILE) \
    go run cmd/app/main.go --keyFile $(KEY_FILE) --domain $(DOMAIN) --key $(KEY) --clientID $(CLIENT_ID) --redirectURI $(REDIRECT_URI)


test:
	go clean -testcache
	go test ./...