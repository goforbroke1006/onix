TRUNCATE service;
TRUNCATE source RESTART IDENTITY;
TRUNCATE release;
TRUNCATE criteria;
TRUNCATE measurement;

INSERT INTO service (title)
VALUES ('foo/bar/backend');

INSERT INTO source (title, kind, address)
VALUES ('stub prometheus', 'prometheus', 'http://127.0.0.1:19091');

INSERT INTO release (service, name, start_at)
VALUES ('foo/bar/backend', '2.1.0', '2020-12-26 00:00:00'), -- 1608940800
       ('foo/bar/backend', '2.0.0', '2020-12-13 00:00:00'), -- 1607817600
       ('foo/bar/backend', '1.2.1', '2020-11-28 00:00:00'),
       ('foo/bar/backend', '1.2.0', '2020-11-14 00:00:00'),
       ('foo/bar/backend', '1.1.0', '2020-11-06 00:00:00'),
       ('foo/bar/backend', '1.0.1', '2020-10-26 00:00:00'),
       ('foo/bar/backend', '1.0.0', '2020-10-25 00:00:00')
;

INSERT INTO criteria (service, title, selector, expected_dir, grouping_interval)
VALUES ('foo/bar/backend', 'processing duration instrument=ONE',
        'histogram_quantile(0.95, sum(increase(api_request_count{environment="prod",instrument="one"}[15m])) by (le))',
        'decrease', '15m'),
       ('foo/bar/backend', 'processing duration instrument=TWO',
        'histogram_quantile(0.95, sum(increase(api_request_count{environment="prod",instrument="two"}[15m])) by (le))',
        'decrease', '15m')
;

INSERT INTO measurement  (source_id, criteria_id, moment, value)
SELECT 1, 1, to_timestamp(_gen_moment), random() * 200 + 100
FROM generate_series(1607817600, 1608940800, 300) AS t(_gen_moment);