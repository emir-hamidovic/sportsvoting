#!/bin/sh

set -e

host="db"
port="3306"
timeout=60

echo "Waiting for the database to become healthy..."

until nc -z -v -w1 "$host" "$port"
do
    if [ "$timeout" -le 0]; then
        echo "Timeout: Database did not become healthy in time"
        exit 1
    fi
    timeout=$((timeout-1))
    sleep 1
done

echo "Database is healthy. Starting the migrate service."
exec "$@"