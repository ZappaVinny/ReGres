-- name: CreateSession :one
INSERT INTO sessions (
    user_id,
    session_token_hash,
    user_agent,
    ip_address,
    expires_at
) VALUES (
    $1, $2, $3, $4, $5
)
RETURNING id, user_id, session_token_hash, user_agent, ip_address, expires_at, created_at, last_seen_at;

-- name: GetSessionByTokenHash :one
SELECT
    sessions.id,
    sessions.user_id,
    sessions.session_token_hash,
    sessions.user_agent,
    sessions.ip_address,
    sessions.expires_at,
    sessions.created_at,
    sessions.last_seen_at,
    users.id AS user_id,
    users.email,
    users.username,
    users.password_hash,
    users.created_at AS user_created_at,
    users.updated_at AS user_updated_at
FROM sessions
JOIN users ON users.id = sessions.user_id
WHERE sessions.session_token_hash = $1
  AND sessions.expires_at > NOW();

-- name: UpdateSessionLastSeen :exec
UPDATE sessions
SET last_seen_at = NOW()
WHERE id = $1;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE session_token_hash = $1;

-- name: DeleteExpiredSessions :exec
DELETE FROM sessions
WHERE expires_at <= NOW();

-- name: DeleteUserSessions :exec
DELETE FROM sessions
WHERE user_id = $1;