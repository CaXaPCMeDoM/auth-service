CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS session
(
    id                 UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    refresh_token_hash VARCHAR(255) UNIQUE NOT NULL,
    jwt_id             uuid UNIQUE         NOT NULL,
    user_ip            INET                NOT NULL,
    expired_at         timestamptz         NOT NULL,
    created_at         timestamptz         NOT NULL,

    user_id            uuid                NOT NULL REFERENCES users (id)
);