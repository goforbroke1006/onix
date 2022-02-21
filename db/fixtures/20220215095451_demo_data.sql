-- migrate:up

INSERT INTO service (title)
VALUES ('foo/backend'),
       ('bar/backend'),
       ('acme/backend');

INSERT INTO source (title, kind, address)
VALUES ('stub prometheus', 'prometheus', 'http://stub-prometheus:19090');

INSERT INTO release (service, name, start_at)
VALUES ('foo/backend', '2.1.0', '2020-12-26 00:00:00'),   -- 1608940800
       ('foo/backend', '2.0.0', '2020-12-13 00:00:00'),   -- 1607817600
       ('foo/backend', '1.2.1', '2020-11-28 00:00:00'),
       ('foo/backend', '1.2.0', '2020-11-14 00:00:00'),
       ('foo/backend', '1.1.0', '2020-11-06 00:00:00'),
       ('foo/backend', '1.0.1', '2020-10-26 00:00:00'),
       ('foo/backend', '1.0.0', '2020-10-25 00:00:00'),

       ('bar/backend', '1.143.0', '2022-02-10 00:00:00'), -- 1644451200
       ('bar/backend', '1.142.0', '2022-02-04 00:00:00'),
       ('bar/backend', '1.141.1', '2022-01-18 00:00:00'),
       ('bar/backend', '1.140.0', '2021-12-12 00:00:00'),
       ('bar/backend', '1.139.0', '2021-12-06 00:00:00'),
       ('bar/backend', '1.138.0', '2021-10-04 00:00:00'),
       ('bar/backend', '1.137.0', '2021-09-30 00:00:00'), -- 1632956400

       ('acme/backend', 'v2.7.0', '2022-02-10 00:00:00'), -- 1644451200
       ('acme/backend', 'v2.6.0', '2022-02-04 00:00:00'),
       ('acme/backend', 'v2.5.1', '2022-01-18 00:00:00'),
       ('acme/backend', 'v2.4.0', '2021-12-12 00:00:00') -- 1639267200
;

INSERT INTO criteria (service, title, selector, expected_dir, grouping_interval)
VALUES ('foo/backend', 'processing duration instrument=ONE',
        'histogram_quantile(0.95, sum(increase(api_request_count{environment="prod",instrument="one"}[15m])) by (le))',
        'decrease', '15m'),
       ('foo/backend', 'processing duration instrument=TWO',
        'histogram_quantile(0.95, sum(increase(api_request_count{environment="prod",instrument="two"}[15m])) by (le))',
        'decrease', '15m')
;

-- fake data for processing duration instrument=ONE
INSERT INTO measurement (source_id, criteria_id, moment, value)
SELECT 1, 1, to_timestamp(_gen_moment), random() * 200 + 100
FROM generate_series(1607817600, 1608940800 + 2 * 12 * 31 * 24 * 60 * 60, 300) AS t(_gen_moment)
ON CONFLICT (source_id, criteria_id, moment) DO UPDATE SET value = EXCLUDED.value;

-- fake data for processing duration instrument=TWO
INSERT INTO measurement (source_id, criteria_id, moment, value)
SELECT 1, 2, to_timestamp(_gen_moment), random() * 200 + 100
FROM generate_series(1607817600, 1608940800 + 2 * 12 * 31 * 24 * 60 * 60, 300) AS t(_gen_moment)
ON CONFLICT (source_id, criteria_id, moment) DO UPDATE SET value = EXCLUDED.value;

-- migrate:down
TRUNCATE service;
TRUNCATE source;
TRUNCATE release;
TRUNCATE criteria;
TRUNCATE measurement;
