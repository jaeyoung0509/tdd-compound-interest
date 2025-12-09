CREATE TABLE users (
    id         CHAR(26)     PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

COMMENT ON TABLE users IS 'User aggregate storing identity info';
COMMENT ON COLUMN users.id IS 'ULID primary key';
