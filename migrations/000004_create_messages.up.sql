CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL REFERENCES users(id),
    client_message_id UUID NOT NULL,
    type VARCHAR(20) NOT NULL,
    content TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,

    CONSTRAINT messages_type_check
        CHECK (type IN ('text', 'image', 'file')),

    CONSTRAINT messages_sender_client_message_unique
        UNIQUE (sender_id, client_message_id)
);

CREATE INDEX messages_conversation_cursor_idx
    ON messages (conversation_id, created_at DESC, id DESC);

ALTER TABLE conversation_members
    ADD COLUMN last_read_message_id UUID REFERENCES messages(id);

CREATE TABLE attachments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    message_id UUID NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    object_key TEXT NOT NULL,
    file_name TEXT NOT NULL,
    content_type VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT attachments_file_size_check
        CHECK (file_size >= 0)
);

CREATE INDEX attachments_message_idx
    ON attachments (message_id);

CREATE TABLE message_receipts (
    message_id UUID NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    delivered_at TIMESTAMPTZ,
    read_at TIMESTAMPTZ,

    PRIMARY KEY (message_id, user_id)
);