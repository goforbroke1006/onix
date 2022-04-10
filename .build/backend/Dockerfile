FROM golang:1.17-alpine AS builder

WORKDIR /code/

RUN go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.9.1
RUN go install github.com/golang/mock/mockgen@v1.6.0

COPY ./go.mod /code/go.mod
COPY ./go.sum /code/go.sum
RUN go mod download

COPY ./ /code/
RUN ls -lah ./
RUN go generate ./...
RUN CGO_ENABLED=0 go build -o onix ./


FROM debian:buster

COPY --from=builder /code/onix /onix

ENTRYPOINT [ "/onix" ]

EXPOSE 8080
