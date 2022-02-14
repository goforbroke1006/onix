TRUNCATE release;

INSERT INTO release (service, name, start_at)
VALUES ('foo/bar/backend', '2.1.0', '2020-12-26 00:00:00'),
       ('foo/bar/backend', '2.0.0', '2020-12-13 00:00:00'),
       ('foo/bar/backend', '1.2.1', '2020-11-28 00:00:00'),
       ('foo/bar/backend', '1.2.0', '2020-11-14 00:00:00'),
       ('foo/bar/backend', '1.1.0', '2020-11-06 00:00:00'),
       ('foo/bar/backend', '1.0.1', '2020-10-26 00:00:00'),
       ('foo/bar/backend', '1.0.0', '2020-10-25 00:00:00')
;