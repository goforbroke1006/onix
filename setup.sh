#!/bin/bash

sudo curl -fsSL -o /usr/local/bin/dbmate https://github.com/amacneil/dbmate/releases/latest/download/dbmate-linux-amd64
sudo chmod +x /usr/local/bin/dbmate

go install github.com/golang/mock/mockgen@v1.6.0
go install golang.org/x/lint/golint
go install github.com/gordonklaus/ineffassign

sudo apt-get install -y gcc g++ make
curl -fsSL https://deb.nodesource.com/setup_16.x | sudo -E bash -
sudo apt install -y nodejs
node --version

sudo chown -R "$(whoami)" "${HOME}/.npm"
npm config set prefix "${HOME}/.npm"
echo 'export PATH="${PATH}:${HOME}/.npm/bin"' >>"${HOME}/.profile"
source "${HOME}/.profile"
npm install -g npm
npm install -g eslint

(
  cd ./web/app/ || exit 1
  npx browserslist@latest --update-db -y
)
