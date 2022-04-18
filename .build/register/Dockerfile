FROM alpine:3.15

ENV ONIX_REGISTER_BASE_URL = "http://127.0.0.1"
ENV ONIX_REGISTER_SERVICE_NAME = ""
ENV ONIX_REGISTER_RELEASE_NAME = ""

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

RUN apk update
RUN apk --no-cache add curl
RUN apk add --no-cache bash

ENTRYPOINT [ "/entrypoint.sh" ]
