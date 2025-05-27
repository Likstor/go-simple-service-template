#!/bin/sh

while IFS= read -r line || [ -n "$line" ]; do
    line=$(echo "$line" | sed -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//' -e 's/#.*//')
    
    [ -z "$line" ] && continue
    
    name="${line%%=*}"
    value="${line#*=}"
    
    name=$(echo "$name" | sed 's/[[:space:]]*$//')
    value=$(echo "$value" | sed 's/^[[:space:]]*//')
    
    [ -n "$name" ] && export "$name=$value"
done < "../configs/.env"

docker compose -f "./docker-compose.yaml" up -d --build