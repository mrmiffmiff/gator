-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, url, feed_id, published_at, title, description)
VALUES (
    $1,
    $2,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT *
FROM posts
WHERE posts.feed_id IN (
    SELECT feed_follows.feed_id
    FROM feed_follows
    WHERE feed_follows.user_id = $1
)
ORDER BY published_at DESC
LIMIT $2;