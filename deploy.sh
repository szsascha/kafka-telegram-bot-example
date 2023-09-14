#!/bin/bash

echo "TELEGRAM_BOT_API_KEY=$TELEGRAM_BOT_API_KEY" > .env

docker compose build

for img in $(docker-compose config | awk '{if ($1 == "image:") print $2;}'); do
  images="$images $img"
done

echo $images

docker image save $images
docker-compose -p "kafka-telegram-bot-example" -H "$DOCKER_REMOTE_HOST" down --rmi all
docker -H "$DOCKER_REMOTE_HOST" image load
docker-compose -p "kafka-telegram-bot-example" -H "$DOCKER_REMOTE_HOST" up --force-recreate -d