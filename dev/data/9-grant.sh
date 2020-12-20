#!/bin/sh
set -e
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "sked" <<-EOSQL

GRANT ALL PRIVILEGES ON DATABASE sked to sked;

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO sked;

EOSQL
