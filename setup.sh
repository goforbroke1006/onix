#!/bin/bash

sudo curl -fsSL -o /usr/local/bin/dbmate https://github.com/amacneil/dbmate/releases/latest/download/dbmate-linux-amd64
sudo chmod +x /usr/local/bin/dbmate

go get -u github.com/deepmap/oapi-codegen/cmd/oapi-codegen

sudo apt-get install -y gcc g++ make
curl -fsSL https://deb.nodesource.com/setup_16.x | sudo -E bash -
sudo apt install -y nodejs
node --version

sudo npm install -g npm

if [[ ! -f config/onix.json ]]; then
  cp config/onix.dist.json config/onix.json
fi
