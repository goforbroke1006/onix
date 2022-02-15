-- migrate:up

INSERT INTO service (title)
VALUES ('foo/bar/backend');

INSERT INTO source (title, kind, address)
VALUES ('stub prometheus', 'prometheus', 'http://stub_prometheus:19090');

INSERT INTO release (service, name, start_at)
VALUES ('foo/bar/backend', '2.1.0', '2020-12-26 00:00:00'), -- 1608940800
       ('foo/bar/backend', '2.0.0', '2020-12-13 00:00:00'), -- 1607817600
       ('foo/bar/backend', '1.2.1', '2020-11-28 00:00:00'),
       ('foo/bar/backend', '1.2.0', '2020-11-14 00:00:00'),
       ('foo/bar/backend', '1.1.0', '2020-11-06 00:00:00'),
       ('foo/bar/backend', '1.0.1', '2020-10-26 00:00:00'),
       ('foo/bar/backend', '1.0.0', '2020-10-25 00:00:00')
;

INSERT INTO criteria (service, title, selector, expected_dir, pull_period)
VALUES ('foo/bar/backend', 'processing duration instrument=ONE',
        'histogram_quantile(0.95, sum(increase(api_request_count{environment="prod",instrument="one"}[15m])) by (le))',
        'decrease', '15m'),
       ('foo/bar/backend', 'processing duration instrument=TWO',
        'histogram_quantile(0.95, sum(increase(api_request_count{environment="prod",instrument="two"}[15m])) by (le))',
        'decrease', '15m')
;

-- fake data for processing duration instrument=ONE
INSERT INTO measurement (source_id, criteria_id, moment, value)
SELECT 1, 1, to_timestamp(_gen_moment), random() * 200 + 100
FROM generate_series(1607817600, 1608940800 + 31 * 24 * 60 * 60, 300) AS t(_gen_moment)
ON CONFLICT (source_id, criteria_id, moment) DO UPDATE SET value = EXCLUDED.value;

-- fake data for processing duration instrument=TWO
INSERT INTO measurement (source_id, criteria_id, moment, value)
SELECT 1, 2, to_timestamp(_gen_moment), random() * 200 + 100
FROM generate_series(1607817600, 1608940800 + 31 * 24 * 60 * 60, 300) AS t(_gen_moment)
ON CONFLICT (source_id, criteria_id, moment) DO UPDATE SET value = EXCLUDED.value;

-- migrate:down
TRUNCATE service;
TRUNCATE source;
TRUNCATE release;
TRUNCATE criteria;
TRUNCATE measurement;
