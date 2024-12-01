FROM golang:1.22.3 AS build

WORKDIR /build

COPY . .

CMD ["go", "test", "./tests/auth_test.go"]
