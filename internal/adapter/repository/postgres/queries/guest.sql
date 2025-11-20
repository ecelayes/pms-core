-- name: CreateGuest :one
INSERT INTO guests (
    hotel_id, email, full_name, phone
) VALUES (
    $1, $2, $3, $4
)
RETURNING id;
