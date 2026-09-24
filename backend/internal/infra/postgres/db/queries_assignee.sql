-- name: ListAssignees :many
SELECT id, name, created_at, updated_at
FROM assignees
ORDER BY name ASC;

-- name: FindAssigneeByName :one
SELECT id, name, created_at, updated_at
FROM assignees
WHERE LOWER(name) = LOWER(sqlc.arg(name)::text)
LIMIT 1;

-- name: FindLeastLoadedAssignee :one
SELECT a.id, a.name, a.created_at, a.updated_at
FROM assignees a
LEFT JOIN tickets t ON t.assignee_id = a.id AND t.status IN ('open', 'in_progress')
GROUP BY a.id, a.name, a.created_at, a.updated_at
ORDER BY COUNT(t.id) ASC, a.name ASC
LIMIT 1;

-- name: UpsertAssignee :exec
INSERT INTO assignees (id, name, created_at, updated_at)
VALUES ($1, $2, $3, $4)
ON CONFLICT (id) DO UPDATE
SET name = EXCLUDED.name,
    updated_at = EXCLUDED.updated_at;
