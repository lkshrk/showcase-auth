
build:
	go build -v ./...

generate-mocks:
	go install go.uber.org/mock/mockgen@latest
	mkdir -p ./pkg/mocks
	go generate ./...

lint:
	go install github.com/mgechev/revive@latest
	revive ./...
	

test:
	go test -v ./...

