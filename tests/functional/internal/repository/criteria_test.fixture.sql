TRUNCATE criteria;

INSERT INTO criteria (service, title, selector, expected_dir)
VALUES ('foo/bar/backend', 'api_registration_rps', 'sum(rate(api_registration_requests_count))', 'equal'),
       ('foo/bar/backend', 'api_login_rps', 'sum(rate(api_login_requests_count))', 'equal'),
       ('foo/bar/backend', 'api_profile_rps', 'sum(rate(api_profile_requests_count))', 'equal')
;