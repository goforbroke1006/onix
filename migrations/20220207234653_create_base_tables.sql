-- migrate:up

SET TIME ZONE 'UTC';

CREATE TABLE service
(
    title VARCHAR(512) PRIMARY KEY
);

CREATE TABLE release
(
    id       SERIAL PRIMARY KEY,
    service  VARCHAR(512) NOT NULL,
    name     VARCHAR(512) NOT NULL,
    start_at TIMESTAMP DEFAULT NOW()
);

ALTER TABLE release
    ADD CONSTRAINT release_service_name_key UNIQUE (service, name);


CREATE TYPE source_type AS ENUM ('prometheus', 'influxdb');

CREATE TABLE source
(
    id      SERIAL PRIMARY KEY,
    title   VARCHAR(1024) NOT NULL UNIQUE,
    kind    SOURCE_TYPE   NOT NULL,
    address VARCHAR(2048) NOT NULL
);

CREATE TYPE dynamic_direction_type AS ENUM ('increase', 'decrease', 'equal');

CREATE TYPE pull_period_type AS ENUM ('30s', '1m', '2m', '5m', '15m');

CREATE TABLE criteria
(
    id           SERIAL PRIMARY KEY,
    service      VARCHAR(512)           NOT NULL,
    title        VARCHAR(512)           NOT NULL,
    selector     TEXT                   NOT NULL,
    expected_dir DYNAMIC_DIRECTION_TYPE NOT NULL,
    pull_period  PULL_PERIOD_TYPE DEFAULT '1m'

);

ALTER TABLE criteria
    ADD CONSTRAINT criteria_service_title_key UNIQUE (service, title);

CREATE TABLE measurement
(
    id          SERIAL PRIMARY KEY,
    source_id   INTEGER         NOT NULL,
    criteria_id INTEGER         NOT NULL,
    moment      TIMESTAMP       NOT NULL,
    value       DECIMAL(24, 12) NOT NULL,
    updated_at  TIMESTAMP DEFAULT NOW()
);

ALTER TABLE measurement
    ADD CONSTRAINT measurement_source_id_criteria_id_moment_key UNIQUE (source_id, criteria_id, moment);

-- migrate:down

