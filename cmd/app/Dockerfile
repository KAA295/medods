FROM golang:1.22.3 AS build

WORKDIR /build
COPY . .

RUN go build -o cmd/app/bin/medods cmd/app/main.go

FROM ubuntu:22.04

WORKDIR app

COPY --from=build /build/cmd/app/bin/medods ./medods

CMD ["sh", "-c", "./medods"]
