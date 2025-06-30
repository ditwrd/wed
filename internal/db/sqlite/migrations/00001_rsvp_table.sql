-- +goose Up
-- +goose StatementBegin
CREATE TABLE rsvps (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    attending BOOLEAN NOT NULL,
    message TEXT,
    group_name TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE rsvps;
-- +goose StatementEnd