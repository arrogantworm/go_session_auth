CREATE TABLE sessions (
    session_id UUID PRIMARY KEY,
    user_id INT REFERENCES users(id) NOT NULL,
    expires_at TIMESTAMP NOT NULL
);