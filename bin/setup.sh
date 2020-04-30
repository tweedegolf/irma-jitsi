#!/bin/bash


export USER_ID="$(id -u)"
export GROUP_ID="$(id -g)"

docker-compose run --no-deps --workdir /go/src/app backend go mod vendor
docker-compose run --no-deps frontend yarn