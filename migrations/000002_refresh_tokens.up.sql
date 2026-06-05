CREATE TABLE advertapp.refresh_tokens(
    token_hash VARCHAR(100) PRIMARY KEY,
    user_id    INTEGER      NOT NULL REFERENCES advertapp.users(id) ON DELETE CASCADE,
    issued_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ  NOT NULL,

    CONSTRAINT expires_after_issued CHECK(expires_at > issued_at) 
);

CREATE INDEX idx_refresh_tokens_user_id ON advertapp.refresh_tokens(user_id);