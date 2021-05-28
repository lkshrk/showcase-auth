
build:
	go build -v ./...

generate-mocks:
	go get github.com/golang/mock/mockgen
	mkdir -p ./pkg/mocks
	mockgen -destination=pkg/mocks/mock_userRepository.go -package=mocks harke.me/showcase-auth/pkg/api UserRepository

test:
	go test -v ./...

