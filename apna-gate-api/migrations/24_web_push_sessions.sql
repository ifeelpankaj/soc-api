-- +migrate Up
ALTER TABLE device_tokens ADD COLUMN session_version BIGINT;
CREATE UNIQUE INDEX device_tokens_web_token_unique ON device_tokens(token)
WHERE platform = 'web' AND session_version IS NOT NULL;

-- +migrate Down
DROP INDEX device_tokens_web_token_unique;
ALTER TABLE device_tokens DROP COLUMN session_version;
