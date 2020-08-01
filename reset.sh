#!/bin/bash
echo 'Stopping containers...'
docker-compose down
echo 'Containers are stopped. Clear data directories...'
rm -Rf data/postgres
rm -Rf data/sqlpad
echo 'Data directories are cleared. Done.'
