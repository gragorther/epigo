-- +goose Up
-- +goose StatementBegin



CREATE TABLE last_messages(
id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
user_id integer references users,
title VARCHAR(200) NOT NULL,
content text
);

CREATE TABLE groups(
id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
user_id integer references users,
name VARCHAR(200) NOT NULL,
description text
);

CREATE TABLE group_last_messages(
group_id integer references groups,
last_message_id integer references last_messages,
PRIMARY KEY (group_id, last_message_id)
);

CREATE TABLE recipients(
id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
group_id integer references groups,
email varchar(319)
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS recipients;
DROP TABLE IF EXISTS group_last_messages;
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS last_messages;

-- +goose StatementEnd
