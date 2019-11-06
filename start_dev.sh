#!/bin/bash

docker build -t mqttparser-dev . -f Dockerfile.dev
docker run --name=mqttparser --network=supergreencloud_back-tier --rm -it -v $(pwd)/config:/etc/mqttparser -v $(pwd):/app mqttparser-dev
docker rmi mqttparser-dev
