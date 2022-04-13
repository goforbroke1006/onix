#!/bin/bash

if [[ -z $BASE_URL ]]; then
  echo "ERROR: specify BASE_URL env"
  exit 1
fi

if [[ -z $SERVICE_NAME ]]; then
  echo "ERROR: specify SERVICE_NAME env"
  exit 1
fi

if [[ -z $RELEASE_NAME ]]; then
  echo "ERROR: specify RELEASE_NAME env"
  exit 1
fi

curl -X GET "$BASE_URL/api/system/register?service_name=$SERVICE_NAME&release_name=$RELEASE_NAME"
