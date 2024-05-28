#!/bin/sh

timestamp=$(date +%s)
SNAKE_NAME=$(echo $NAME | sed 's/[A-Z]/_\l&/g')
mkdir -p ./migrations/$timestamp$SNAKE_NAME
touch ./migrations/$timestamp$SNAKE_NAME/$timestamp$SNAKE_NAME.up.sql
touch ./migrations/$timestamp$SNAKE_NAME/$timestamp$SNAKE_NAME.down.sql
