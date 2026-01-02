 -- name: CreateInterview :one
INSERT INTO interviews (
    id,
    title,
    message_count,
    prompt_tokens,
    completion_tokens,
    cost,
    summary_message_id,
    todos,
    difficulty,
    topic,
    status,
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
    ?,
    ?,
    ?,
    ?,
    strftime('%s', 'now'),
    strftime('%s', 'now')
) RETURNING *;

-- name: GetInterviewByID :one
SELECT *
FROM interviews
WHERE id = ? LIMIT 1;

-- name: ListInterviews :many
SELECT *
FROM interviews
ORDER BY updated_at DESC;

-- name: UpdateInterview :one
UPDATE interviews
SET
    title = ?,
    message_count = ?,
    prompt_tokens = ?,
    completion_tokens = ?,
    cost = ?,
    summary_message_id = ?,
    todos = ?,
    difficulty = ?,
    topic = ?,
    status = ?
WHERE id = ?
RETURNING id, title, message_count, prompt_tokens, completion_tokens, cost, summary_message_id, todos, difficulty, topic, status, updated_at, created_at;

-- name: UpdateInterviewTitleAndUsage :exec
UPDATE interviews
SET
    title = ?,
    prompt_tokens = prompt_tokens + ?,
    completion_tokens = completion_tokens + ?,
    cost = cost + ?
WHERE id = ?;

-- name: DeleteInterview :exec
DELETE FROM interviews
WHERE id = ?;

