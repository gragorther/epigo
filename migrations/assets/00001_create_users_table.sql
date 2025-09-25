-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    email VARCHAR(319) UNIQUE NOT NULL,
    name VARCHAR(70),
    username VARCHAR(50) UNIQUE NOT NULL,
    last_login TIMESTAMP WITH TIME ZONE,
    cron VARCHAR(20) NOT NULL DEFAULT '0 0 * * 0',
    sent_emails SMALLINT NOT NULL DEFAULT 0,
    max_sent_emails SMALLINT NOT NULL DEFAULT 7,
    admin BOOLEAN DEFAULT FALSE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE FUNCTION set_updated_at_column() RETURNS trigger AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_set_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION set_updated_at_column();
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS users_set_updated_at;
DROP FUNCTION IF EXISTS set_updated_at_column;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
