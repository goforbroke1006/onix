FROM node:16-buster AS builder

WORKDIR /code/

COPY ./frontend/dashboard-admin/package.json      /code/package.json
COPY ./frontend/dashboard-admin/package-lock.json /code/package-lock.json
RUN npm install

COPY ./frontend/dashboard-admin/ /code/
RUN npm run build



FROM nginx:1.21.6

ENV REACT_APP_API_DASHBOARD_ADMIN_BASE_ADDR ''

COPY --from=builder /code/build/ /usr/share/nginx/html/

COPY ./.build/frontend/main/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT [ "/entrypoint.sh" ]
