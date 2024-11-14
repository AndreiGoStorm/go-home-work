-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE events (
    id VARCHAR(36) NOT NULL,
    title VARCHAR(255) NOT NULL,
    start TIMESTAMP NOT NULL,
    finish TIMESTAMP NOT NULL,
    description TEXT NULL,
    user_id VARCHAR(36) NOT NULL,
    remind INT NOT NULL DEFAULT 0
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE events;
