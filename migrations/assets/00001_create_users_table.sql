-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(319) UNIQUE NOT NULL,
    name VARCHAR(70),
    username VARCHAR(50) UNIQUE NOT NULL,
    last_login TIMESTAMP WITH TIME ZONE,
    cron VARCHAR(20),
    sent_emails SMALLINT,
    max_sent_emails SMALLINT,
    admin BOOLEAN DEFAULT FALSE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
