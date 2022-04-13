FROM alpine:3.15

ENV BASE_URL = "http://127.0.0.1"
ENV SERVICE_NAME = ""
ENV RELEASE_NAME = ""

COPY entrypoint.sh /entrypoint.sh
RUN apk --no-cache add curl

ENTRYPOINT [ "/entrypoint.sh" ]
