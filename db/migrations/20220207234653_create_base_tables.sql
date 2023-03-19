-- migrate:up

SET TIME ZONE 'UTC';

CREATE TABLE service
(
    id VARCHAR(512) PRIMARY KEY
);

CREATE TABLE release
(
    service  VARCHAR(512) NOT NULL,
    tag      VARCHAR(512) NOT NULL,
    start_at TIMESTAMP DEFAULT NOW()
);

ALTER TABLE release
    ADD CONSTRAINT release_service_name_key UNIQUE (service, tag);


CREATE TYPE source_type AS ENUM ('prometheus', 'influxdb');

CREATE TABLE source
(
    id      VARCHAR(1024) NOT NULL PRIMARY KEY,
    kind    SOURCE_TYPE   NOT NULL,
    address VARCHAR(2048) NOT NULL
);

CREATE TYPE dynamic_direction_type AS ENUM ('increase', 'decrease', 'equal');

CREATE TYPE grouping_interval_type AS ENUM ('30s', '1m', '2m', '5m', '15m');

CREATE TABLE criteria
(
    id        SERIAL PRIMARY KEY,
    service   VARCHAR(512)           NOT NULL,
    title     VARCHAR(512)           NOT NULL,
    selector  TEXT                   NOT NULL,
    direction DYNAMIC_DIRECTION_TYPE NOT NULL,
    interval  GROUPING_INTERVAL_TYPE DEFAULT '1m'
);

ALTER TABLE criteria
    ADD CONSTRAINT criteria_service_title_key UNIQUE (service, title);

CREATE TABLE IF NOT EXISTS measurement
(
    criteria_id INT NOT NULL,
    moment      BIGINT,
    value       DECIMAL(24, 6),
    updated_at  TIMESTAMP DEFAULT NOW()
);

ALTER TABLE measurement
    ADD CONSTRAINT measurement_criteria_ts_key UNIQUE (criteria_id, moment);

CREATE INDEX IF NOT EXISTS measurement_criteria_ts_idx ON measurement (criteria_id, moment);


-- migrate:down

