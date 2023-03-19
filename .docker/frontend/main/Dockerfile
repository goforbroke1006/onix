FROM node:16-buster AS builder

WORKDIR /code/

COPY ./web/dashboard-main-app/package.json      /code/package.json
COPY ./web/dashboard-main-app/package-lock.json /code/package-lock.json
RUN npm install

COPY ./web/dashboard-main-app/ /code/
RUN npm run build



FROM nginx:1.21.6

COPY --from=builder /code/build/ /usr/share/nginx/html/

COPY ./.build/frontend/main/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
ENTRYPOINT [ "/entrypoint.sh" ]
