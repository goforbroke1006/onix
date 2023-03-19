FROM golang:1.19-alpine AS builder

WORKDIR /code/

COPY ./go.mod /code/go.mod
COPY ./go.sum /code/go.sum
RUN go mod download

COPY ./ /code/
RUN go generate ./...
RUN CGO_ENABLED=0 go build -o onix ./


FROM debian:buster

COPY --from=builder /code/onix /onix
COPY ./db/migrations /db/migrations

ENTRYPOINT [ "/onix" ]

EXPOSE 8080
