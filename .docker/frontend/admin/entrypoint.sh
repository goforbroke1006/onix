#!/bin/bash

REACT_APP_API_DASHBOARD_ADMIN_BASE_ADDR=${REACT_APP_API_DASHBOARD_ADMIN_BASE_ADDR:-'http://127.0.0.1:8080/api/dashboard-admin'}

if [[ -n ${REACT_APP_API_DASHBOARD_ADMIN_BASE_ADDR} ]]; then
  echo 'window.env = {"REACT_APP_API_DASHBOARD_ADMIN_BASE_ADDR": "'"${REACT_APP_API_DASHBOARD_ADMIN_BASE_ADDR}"'"}' \
    >/usr/share/nginx/html/env.js
fi


nginx -g 'daemon off;'
