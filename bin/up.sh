#!/bin/bash

export USER_ID="$(id -u)"
export GROUP_ID="$(id -g)"
HOST_LAN_IP=$(hostname -I | awk '{print $1}') docker-compose up