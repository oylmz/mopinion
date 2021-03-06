export GO111MODULE=on
export GOFLAGS=-mod=vendor
export CGO_ENABLED=0

build: compile

compile: 
	go build ./...
	go vet ./...

lint:
	# golint doesn't exclude vendor from ./... yet
	golint $$(GOFLAGS=-mod=vendor go list ./...)

test: lint
	CGO_ENABLED=1 go test -race -v -p 1 -timeout 60m -covermode=atomic -coverprofile coverage.txt ./...
	go tool cover -func=coverage.txt | grep "total"
	go tool cover -html=coverage.txt -o coverage.html

fmt:
	go fmt ./...
