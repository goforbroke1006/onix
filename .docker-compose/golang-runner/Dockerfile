FROM golang:1.16-buster

COPY .docker-compose/golang-runner/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT [ "bash", "/entrypoint.sh" ]

EXPOSE 8080
