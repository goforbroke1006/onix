#!/bin/bash

sudo curl -fsSL -o /usr/local/bin/dbmate https://github.com/amacneil/dbmate/releases/latest/download/dbmate-linux-amd64
sudo chmod +x /usr/local/bin/dbmate

go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.44.2
go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen
go install github.com/golang/mock/mockgen@v1.6.0
go install golang.org/x/lint/golint
go get -u github.com/gordonklaus/ineffassign
go install github.com/client9/misspell/cmd/misspell

sudo apt-get install -y gcc g++ make
curl -fsSL https://deb.nodesource.com/setup_16.x | sudo -E bash -
sudo apt install -y nodejs
node --version

sudo chown -R "$(whoami)" "${HOME}/.npm"
npm config set prefix "${HOME}/.npm"
echo 'export PATH="${PATH}:${HOME}/.npm/bin"' >> "${HOME}/.profile"
source "${HOME}/.profile"
npm install -g npm
npm install -g eslint

if [[ ! -f config/onix.json ]]; then
  cp config/onix.dist.json config/onix.json
fi
