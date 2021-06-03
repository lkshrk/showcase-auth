
build:
	go build -v ./...

generate-mocks:
	go get github.com/golang/mock/mockgen
	mkdir -p ./pkg/mocks
	go generate ./...

test:
	go test -v ./...

