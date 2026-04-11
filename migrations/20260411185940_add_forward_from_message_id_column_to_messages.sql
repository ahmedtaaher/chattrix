-- +goose Up
-- +goose StatementBegin
ALTER TABLE messages ADD COLUMN forward_from_message_id UUID REFERENCES messages(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE messages DROP COLUMN forward_from_message_id;
-- +goose StatementEnd
