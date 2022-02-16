FROM node:16-buster

WORKDIR /code/

COPY .docker-compose/frontend-prod-builder/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT [ "bash", "/entrypoint.sh" ]

HEALTHCHECK --timeout=10s --interval=5s --retries=10 CMD "bash -c '[ -f /code/build/index.html ]'"
