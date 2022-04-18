#!/bin/bash

if [[ -z ${ONIX_REGISTER_BASE_URL} ]]; then
  echo "ERROR: specify ONIX_REGISTER_BASE_URL env"
  exit 1
fi

if [[ -z ${ONIX_REGISTER_SERVICE_NAME} ]]; then
  echo "ERROR: specify ONIX_REGISTER_SERVICE_NAME env"
  exit 1
fi

if [[ -z ${ONIX_REGISTER_RELEASE_NAME} ]]; then
  echo "ERROR: specify ONIX_REGISTER_RELEASE_NAME env"
  exit 1
fi

url="${ONIX_REGISTER_BASE_URL}/api/system/register?service_name=${ONIX_REGISTER_SERVICE_NAME}&release_name=${ONIX_REGISTER_RELEASE_NAME}"
echo "${url}"
curl -X GET "${url}"
