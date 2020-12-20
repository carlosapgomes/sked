#!/bin/sh
set -e
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL

CREATE USER sked WITH ENCRYPTED PASSWORD 'sked';

CREATE DATABASE sked;

EOSQL
