-- set of queries to populate a db locally for development purposes

START TRANSACTION;

INSERT INTO calendar (date, event)
VALUES
    ('10 y 11 de julio', 'finde en la costa'),
    ('13 de julio 9:00', 'dentista'),
    ('18 de julio', 'cumple Mengano'),
    ('24 de julio', 'cumple Fulano');

COMMIT;
