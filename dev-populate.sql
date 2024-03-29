-- set of queries to populate a db locally for development purposes

START TRANSACTION;

INSERT INTO calendar (date, event)
VALUES
    ('18 de julio', 'cumple Mengano'),
    ('24 de julio', 'cumple Fulano');

COMMIT;
