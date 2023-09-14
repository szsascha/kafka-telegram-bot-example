#!/bin/bash

docker compose build

for img in $(docker-compose config | awk '{if ($1 == "image:") print $2;}'); do
  images="$images $img"
done

echo $images

docker image save $images
docker-compose -p "kafka-telegram-bot-example" -H "$DOCKER_REMOTE_HOST" down --rmi all
docker -H "$DOCKER_REMOTE_HOST" image load
echo "TELEGRAM_BOT_API_KEY=$TELEGRAM_BOT_API_KEY" > .actions_env
docker-compose -p "kafka-telegram-bot-example" -H "$DOCKER_REMOTE_HOST" --env-file ./.actions_env up --force-recreate -d