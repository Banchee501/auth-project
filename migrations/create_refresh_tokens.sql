CREATE TABLE refresh_tokens (
  id SERIAL PRIMARY KEY,
  user_id INT REFERENCES users(id) ON DELETE CASCADE,
  token TEXT NOT NULL,
  device_id TEXT,
  user_agent TEXT,
  ip TEXT,
  expires_at TIMESTAMP NOT NULL
);