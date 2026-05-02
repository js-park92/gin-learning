CREATE TABLE users (
    id         SERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
    name       VARCHAR(255) NOT NULL,
    email      VARCHAR(255) NOT NULL
);

CREATE UNIQUE INDEX idx_users_email     ON users (email) WHERE deleted_at IS NULL;
CREATE        INDEX idx_users_deleted_at ON users (deleted_at);
