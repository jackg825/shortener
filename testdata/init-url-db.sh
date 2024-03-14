#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE TABLE urls (
        short_url VARCHAR(255) PRIMARY KEY,
        long_url TEXT NOT NULL,
        user_id VARCHAR(40)
    );
EOSQL

