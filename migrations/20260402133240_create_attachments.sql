-- +goose Up
-- +goose StatementBegin
CREATE TABLE attachments (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  message_id UUID NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
  file_url TEXT NOT NULL,
  file_type VARCHAR(50),
  file_size BIGINT
);
CREATE INDEX idx_attachments_message_id ON attachments(message_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE attachments;
-- +goose StatementEnd
