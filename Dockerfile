FROM golang:alpine AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o main cmd/server/main.go

WORKDIR /dist

RUN cp /build/main .


FROM scratch
COPY --from=builder /dist/main /
ENTRYPOINT ["/main"]