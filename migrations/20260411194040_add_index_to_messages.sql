-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_messages_chat_id_sent_at ON messages(chat_id, sent_at DESC);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX idx_messages_chat_id_sent_at;
-- +goose StatementEnd
