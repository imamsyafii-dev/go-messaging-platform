DROP TABLE IF EXISTS message_receipts;
DROP TABLE IF EXISTS attachments;

ALTER TABLE conversation_members
    DROP COLUMN IF EXISTS last_read_message_id;

DROP TABLE IF EXISTS messages;