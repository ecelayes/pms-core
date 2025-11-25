-- name: CreateUser :one
INSERT INTO users (hotel_id, email, password_hash, role)
VALUES (
    sqlc.narg('hotel_id'),
    $1, $2, $3
)
RETURNING id;

-- name: GetUserByEmail :one
SELECT id, hotel_id, email, password_hash, role 
FROM users 
WHERE email = $1 LIMIT 1;
