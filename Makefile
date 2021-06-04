
build:
	go build -v ./...

generate-mocks:
	go get github.com/golang/mock/mockgen
	mkdir -p ./pkg/mocks
	go generate ./...

lint:
	go get -u golang.org/x/lint/golint
	golint ./...
	

test:
	go test -v ./...

