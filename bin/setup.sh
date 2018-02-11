#!/bin/bash
docker-compose run action-instance-service dep ensure
docker-compose run action-instance-service go run migrate/migrate.go up