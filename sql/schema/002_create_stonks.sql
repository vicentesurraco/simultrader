-- migrations/002_create_stonks.sql
-- +goose Up
CREATE TABLE stonks(
	id SERIAL PRIMARY KEY,
	user_id INTEGER NOT NULL,
	symbol VARCHAR(10) NOT NULL,
	is_active BOOLEAN DEFAULT TRUE,
	UNIQUE (user_id, symbol),
	FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE stonks;
