-- set of queries to populate a db locally for development purposes

START TRANSACTION;

INSERT INTO events (starts_at, description)
VALUES
    ('2024-12-07 12:00:00', 'finde formación'),
    ('2024-07-21 15:00:00', 'cumple Aurora')
;

COMMIT;
