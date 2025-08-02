-- +goose Up
-- +goose StatementBegin
CREATE TABLE tasks (
 id SERIAL PRIMARY KEY,
 webhook TEXT NOT NULL,
 send_date TIMESTAMP NOT NULL,
 payload JSONB NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tasks;
-- +goose StatementEnd
