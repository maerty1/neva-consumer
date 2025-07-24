#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
CREATE TABLE if not exists accounts
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    token TEXT NOT NULL,
    server_host VARCHAR(255) NOT NULL
);

CREATE TABLE if not exists scada_measure_points 
(
    id SERIAL PRIMARY KEY,
    account_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC'),

    FOREIGN KEY (account_id) REFERENCES accounts(id),
    UNIQUE (account_id, title)
);

CREATE TABLE if not exists scada_rawdata 
(
    scada_measure_point_id INTEGER NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    varname TEXT NULL,
    measurement_type_id INTEGER NULL,
    value TEXT,
    raw_packet JSONB
);

INSERT INTO public.accounts (id, name, token, server_host) VALUES (1, 'Test', 'Test', 'Test');

EOSQL