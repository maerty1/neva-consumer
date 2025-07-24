#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL

SET timezone = 'UTC';

CREATE TABLE scada_consumer
(
    data_type VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC'), 
    other_data JSONB
);

CREATE TABLE if not exists accounts
(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    token TEXT NOT NULL,
    server_host VARCHAR(255) NOT NULL
);

CREATE TABLE if not exists measure_points 
(
    id SERIAL PRIMARY KEY,
    account_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    address VARCHAR(255) NULL,
    full_title VARCHAR(255) NULL,
    device_id INTEGER NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC'),

    FOREIGN KEY (account_id) REFERENCES accounts(id)
);

CREATE TABLE if not exists measure_points_data
(   
    measure_point_id INTEGER NOT NULL,
    datetime TIMESTAMP NOT NULL,
    values JSONB,

    PRIMARY KEY (measure_point_id, datetime)
);

CREATE TABLE if not exists measure_points_data_day
(   
    measure_point_id INTEGER NOT NULL,
    datetime TIMESTAMP NOT NULL,
    values JSONB,

    PRIMARY KEY (measure_point_id, datetime)
);

CREATE TABLE if not exists accounts_sync_log
(
    account_id INTEGER NOT NULL,
    measure_point_id INTEGER,
    message TEXT,
    level VARCHAR(255) NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC')
);

CREATE TABLE if not exists measure_points_poll_log 
(
    poll_id INTEGER NOT NULL,
    -- Ответ из Лэрс во время попытки опроса --
    message TEXT,
    -- Результат опроса --
    status TEXT,
    measure_point_id INTEGER,
    account_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC')
);

CREATE TABLE if not exists measure_points_poll_retry 
(
    original_poll_id INTEGER NOT NULL,
    retrying_poll_id INTEGER NOT NULL,
    status TEXT DEFAULT 'PENDING', 
    created_at TIMESTAMP NOT NULL DEFAULT (current_timestamp AT TIME ZONE 'UTC')
);

CREATE INDEX ix_measure_points_data ON public.measure_points_data  USING btree ("datetime", measure_point_id);

EOSQL