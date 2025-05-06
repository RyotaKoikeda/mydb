BEGIN;

DELETE FROM dummy_users WHERE email IN ('alice@example.com', 'bob@example.com');

COMMIT;