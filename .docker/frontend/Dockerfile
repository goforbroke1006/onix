FROM node:16-buster AS builder

WORKDIR /code/

COPY ./web/app/package.json      /code/package.json
COPY ./web/app/package-lock.json /code/package-lock.json
RUN npm install

COPY ./web/app/ /code/
RUN npm run build


FROM nginx:1.21.6

COPY --from=builder /code/build/ /usr/share/nginx/html/

COPY ./.docker/frontend/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENTRYPOINT [ "/entrypoint.sh" ]
