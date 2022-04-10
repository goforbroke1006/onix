FROM node:16-buster AS builder

WORKDIR /code/
COPY ./package.json /code/package.json
COPY ./package-lock.json /code/package-lock.json
RUN npm install

COPY ./ /code/
RUN npm run build


FROM nginx:1.21.6

COPY --from=builder /code/build/ /usr/share/nginx/html/
