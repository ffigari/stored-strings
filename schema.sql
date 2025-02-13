START TRANSACTION;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS events (
    id          UUID PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
    starts_at   timestamp NOT NULL,
    description TEXT NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT NOW(),
    updated_at  timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS clicks (
    id UUID PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
    x  int not null,
    y  int not null
);

CREATE TABLE IF NOT EXISTS generic_ars_expenses (
    id          UUID PRIMARY KEY DEFAULT UUID_GENERATE_V4(),
    description TEXT NOT NULL,
    ars_cents   BIGINT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;
