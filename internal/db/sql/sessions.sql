-- name: CreateSession :one
INSERT INTO sessions (
    id,
    parent_session_id,
    title,
    message_count,
    prompt_tokens,
    completion_tokens,
    cost,
    summary_message_id,
    mode,
    updated_at,
    created_at
) VALUES (
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    ?,
    null,
    ?,
    strftime('%s', 'now'),
    strftime('%s', 'now')
) RETURNING id, parent_session_id, title, message_count, prompt_tokens, completion_tokens, cost, updated_at, created_at, summary_message_id, todos, mode;

-- name: GetSessionByID :one
SELECT id, parent_session_id, title, message_count, prompt_tokens, completion_tokens, cost, updated_at, created_at, summary_message_id, todos, mode
FROM sessions
WHERE id = ? LIMIT 1;

-- name: ListSessions :many
SELECT id, parent_session_id, title, message_count, prompt_tokens, completion_tokens, cost, updated_at, created_at, summary_message_id, todos, mode
FROM sessions
WHERE parent_session_id is NULL
ORDER BY updated_at DESC;

-- name: UpdateSession :one
UPDATE sessions
SET
    title = ?,
    prompt_tokens = ?,
    completion_tokens = ?,
    summary_message_id = ?,
    cost = ?,
    todos = ?,
    mode = ?
WHERE id = ?
RETURNING id, parent_session_id, title, message_count, prompt_tokens, completion_tokens, cost, updated_at, created_at, summary_message_id, todos, mode;

-- name: UpdateSessionTitleAndUsage :exec
UPDATE sessions
SET
    title = ?,
    prompt_tokens = prompt_tokens + ?,
    completion_tokens = completion_tokens + ?,
    cost = cost + ?
WHERE id = ?;


-- name: DeleteSession :exec
DELETE FROM sessions
WHERE id = ?;
