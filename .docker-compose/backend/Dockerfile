FROM golang:1.18-bullseye

RUN apt update && apt upgrade -y
RUN apt install -y curl

RUN go install github.com/githubnemo/CompileDaemon@v1.4.0

#ARG LDAP_USERNAME=''
#ARG LDAP_PASSWORD=''
#RUN echo "machine your.private.gitlab.etc login "$LDAP_USERNAME" password "$LDAP_PASSWORD > /root/.netrc
#
#ENV GOPRIVATE=your.private.gitlab.etc

WORKDIR /code/
COPY ./go.mod /code/go.mod
COPY ./go.sum /code/go.sum
RUN go mod download -x

COPY . /code/
RUN go generate ./...

COPY ./.docker-compose/backend/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENTRYPOINT [ "/entrypoint.sh" ]
