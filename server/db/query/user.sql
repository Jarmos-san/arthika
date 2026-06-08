-- name: CreateUser :exec
INSERT INTO users (id, username, email, password_hash)
VALUES (?, ?, ?, ?);

-- name: FindUserByEmail :one
SELECT id, username, email, password_hash
FROM users
WHERE email = ? LIMIT 1;
