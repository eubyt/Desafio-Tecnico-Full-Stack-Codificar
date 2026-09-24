-- name: CreateTicket :exec
INSERT INTO tickets (id, title, description, priority, status, assignee_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: UpdateTicket :execrows
UPDATE tickets
SET title = $2,
    description = $3,
    priority = $4,
    status = $5,
    assignee_id = $6,
    updated_at = $7
WHERE id = $1;

-- name: FindTicketByID :one
SELECT 
    t.id, t.title, t.description, t.priority, t.status, t.assignee_id, t.created_at, t.updated_at,
    a.id AS assignee_ref_id, a.name AS assignee_name, a.created_at AS assignee_created_at, a.updated_at AS assignee_updated_at
FROM tickets t
JOIN assignees a ON a.id = t.assignee_id
WHERE t.id = $1
LIMIT 1;

-- name: CountTicketsWithFilter :one
SELECT COUNT(*)
FROM tickets t
JOIN assignees a ON a.id = t.assignee_id
WHERE (sqlc.narg('status')::text IS NULL OR t.status = sqlc.narg('status'))
  AND (sqlc.narg('priority')::text IS NULL OR t.priority = sqlc.narg('priority'))
  AND (sqlc.narg('assignee')::text IS NULL OR LOWER(a.name) = LOWER(sqlc.narg('assignee')))
  AND (sqlc.narg('search')::text IS NULL OR LOWER(t.title) LIKE ('%' || LOWER(sqlc.narg('search')::text) || '%') OR LOWER(t.description) LIKE ('%' || LOWER(sqlc.narg('search')::text) || '%'));

-- name: ListTicketsWithFilterAndSort :many
SELECT 
    t.id, t.title, t.description, t.priority, t.status, t.assignee_id, t.created_at, t.updated_at,
    a.id AS assignee_ref_id, a.name AS assignee_name, a.created_at AS assignee_created_at, a.updated_at AS assignee_updated_at
FROM tickets t
JOIN assignees a ON a.id = t.assignee_id
WHERE (sqlc.narg('status')::text IS NULL OR t.status = sqlc.narg('status'))
  AND (sqlc.narg('priority')::text IS NULL OR t.priority = sqlc.narg('priority'))
  AND (sqlc.narg('assignee')::text IS NULL OR LOWER(a.name) = LOWER(sqlc.narg('assignee')))
  AND (sqlc.narg('search')::text IS NULL OR LOWER(t.title) LIKE ('%' || LOWER(sqlc.narg('search')::text) || '%') OR LOWER(t.description) LIKE ('%' || LOWER(sqlc.narg('search')::text) || '%'))
ORDER BY
  -- Priority ASC
  CASE WHEN sqlc.arg('sort_by')::text = 'priority' AND sqlc.arg('sort_order')::text = 'asc' THEN (CASE t.priority WHEN 'low' THEN 1 WHEN 'medium' THEN 2 WHEN 'high' THEN 3 ELSE 0 END) END ASC,
  -- Priority DESC
  CASE WHEN sqlc.arg('sort_by')::text = 'priority' AND sqlc.arg('sort_order')::text = 'desc' THEN (CASE t.priority WHEN 'low' THEN 1 WHEN 'medium' THEN 2 WHEN 'high' THEN 3 ELSE 0 END) END DESC,

  -- Status ASC
  CASE WHEN sqlc.arg('sort_by')::text = 'status' AND sqlc.arg('sort_order')::text = 'asc' THEN t.status END ASC,
  -- Status DESC
  CASE WHEN sqlc.arg('sort_by')::text = 'status' AND sqlc.arg('sort_order')::text = 'desc' THEN t.status END DESC,

  -- Title ASC
  CASE WHEN sqlc.arg('sort_by')::text = 'title' AND sqlc.arg('sort_order')::text = 'asc' THEN t.title END ASC,
  -- Title DESC
  CASE WHEN sqlc.arg('sort_by')::text = 'title' AND sqlc.arg('sort_order')::text = 'desc' THEN t.title END DESC,

  -- CreatedAt ASC
  CASE WHEN sqlc.arg('sort_by')::text = 'created_at' AND sqlc.arg('sort_order')::text = 'asc' THEN t.created_at END ASC,
  -- CreatedAt DESC (default fallback)
  CASE WHEN (sqlc.arg('sort_by')::text = 'created_at' AND sqlc.arg('sort_order')::text = 'desc') OR (sqlc.arg('sort_by')::text NOT IN ('priority', 'status', 'title', 'created_at')) THEN t.created_at END DESC,

  t.id DESC
OFFSET sqlc.arg('page_offset')::int
LIMIT sqlc.arg('page_limit')::int;

-- name: CountOpenTicketsByAssignee :many
SELECT a.name AS assignee_name, COUNT(t.id)::int AS open_count
FROM tickets t
JOIN assignees a ON a.id = t.assignee_id
WHERE t.status IN ('open', 'in_progress')
GROUP BY a.name;
