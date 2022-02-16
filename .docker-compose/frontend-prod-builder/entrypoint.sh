#!/bin/bash

npm install
chmod -R 0777 /code/node_modules/
chmod 0777 /code/package-lock.json

mkdir -p /code/build/
chmod -R 0777 /code/build/
npm run build
