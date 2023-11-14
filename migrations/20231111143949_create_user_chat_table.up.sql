CREATE TABLE IF NOT EXISTS user_chat (
    id BIGSERIAL PRIMARY KEY,
    "from" BIGINT NOT NULL,
    "to" BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL,
    deleted_at TIMESTAMP DEFAULT NULL,
    CONSTRAINT fk_from_user FOREIGN KEY ("from") REFERENCES users(id),
    CONSTRAINT fk_to_user FOREIGN KEY ("to") REFERENCES users(id)

);

CREATE UNIQUE INDEX userchat_id_idx ON user_chat (id);
CREATE INDEX userchat_from_idx ON user_chat ("from");
CREATE INDEX userchat_to_idx ON user_chat ("to");