CREATE TABLE IF NOT EXISTS messages (
    id BIGSERIAL PRIMARY KEY,
    user_chat_id BIGINT NOT NULL,
    message_text TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL,
    deleted_at TIMESTAMP DEFAULT NULL,
    CONSTRAINT fk_messages
      FOREIGN KEY(user_chat_id) 
	    REFERENCES user_chat(id)
);

CREATE UNIQUE INDEX messages_id_idx ON messages (id);
CREATE INDEX messages_user_chat_id_idx ON messages (user_chat_id);
CREATE INDEX messages_message_text_idx ON messages (message_text);