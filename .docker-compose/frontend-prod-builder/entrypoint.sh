#!/bin/bash

npm install
chmod -R 0777 ./node_modules/
chmod 0777 ./package-lock.json

mkdir -p ./build/
chmod -R 0777 ./build/

mkdir -p ./dist/
chmod -R 0777 ./dist/

BUILD_PATH='./dist' npm run build

rm -rf ./build/*
mv ./dist/* ./build/
rm -r ./dist/
