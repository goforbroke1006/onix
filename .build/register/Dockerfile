FROM alpine:3.15

ENV BASE_URL = "http://127.0.0.1"
ENV SERVICE_NAME = ""
ENV RELEASE_NAME = ""

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

RUN apk update
RUN apk --no-cache add curl
RUN apk add --no-cache bash

ENTRYPOINT [ "/entrypoint.sh" ]
