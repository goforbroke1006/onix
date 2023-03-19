FROM golang:1.19-alpine AS builder

WORKDIR /code/

COPY ./go.mod /code/go.mod
COPY ./go.sum /code/go.sum
RUN go mod download

ENV CGO_ENABLED=0
COPY ./ /code/
RUN go generate ./...
RUN go build -o onix ./


FROM debian:buster

COPY --from=builder /code/onix /onix

ENTRYPOINT [ "/onix" ]

EXPOSE 8080
