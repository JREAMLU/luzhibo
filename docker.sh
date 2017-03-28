#!/bin/sh

curl -H "Content-Type: application/json" --data '{"docker_tag": "`$DOCKER_TAG`"}' -X POST https://registry.hub.docker.com/u/$DOCKER_USER/$DOCKER_REPO/trigger/$DOCKER_TRIGGER_TOKEN/
