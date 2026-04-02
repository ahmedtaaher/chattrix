-- +goose Up
-- +goose StatementBegin
CREATE TABLE chat_members (
  chat_id UUID REFERENCES chats(id) ON DELETE CASCADE,
  user_id UUID REFERENCES users(id) ON DELETE CASCADE,
  role VARCHAR(20) DEFAULT 'member',
  joined_at TIMESTAMPTZ DEFAULT NOW(),
  is_pinned BOOLEAN DEFAULT FALSE,
  is_muted BOOLEAN DEFAULT FALSE,
  last_read_message_id UUID,
  PRIMARY KEY (chat_id, user_id)
);
CREATE INDEX idx_chat_members_user_id ON chat_members(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE chat_members;
-- +goose StatementEnd
