#!/bin/sh

/migrator \
    -storage-URL=${STORAGE_USER}:${STORAGE_PASSWORD}@${STORAGE_HOST}:${STORAGE_PORT}/${STORAGE_NAME} \
    -migrations-path=migrations \
    -migrations-table=migrations \
    -action=up

exec ./service