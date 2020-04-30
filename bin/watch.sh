#!/bin/bash

# Watch for file changes an re-run the Go application

go run "$@" &
while inotifywait --exclude .swp -e modify -r . ;
do
    pkill -f "/tmp/go-build.*/b001/exe/$1"
    go run "$@" &
done;
