FROM golang:1.19-bullseye

RUN apt update && apt upgrade -y
RUN apt install -y curl

RUN go install github.com/githubnemo/CompileDaemon@v1.4.0

WORKDIR /code/
COPY ./go.mod /code/go.mod
COPY ./go.sum /code/go.sum
RUN go mod download -x

COPY ./cmd      /code/cmd
COPY ./domain   /code/domain
COPY ./external /code/external
COPY ./internal /code/internal
COPY ./pkg      /code/pkg
COPY ./main.go  /code/main.go
RUN go generate ./...

COPY ./.docker-compose/backend/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENTRYPOINT [ "/entrypoint.sh" ]

EXPOSE 8080
