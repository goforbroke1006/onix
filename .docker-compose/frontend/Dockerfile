FROM node:16-buster AS builder

WORKDIR /code/

COPY ./web/app/package.json      /code/package.json
COPY ./web/app/package-lock.json /code/package-lock.json
RUN npm install

COPY ./web/app/ /code/
ENTRYPOINT npm run start

EXPOSE 3000
