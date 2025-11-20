-- name: CreateGuest :one
INSERT INTO guests (
    hotel_id, email, full_name, phone
) VALUES (
    $1, $2, $3, $4
)
ON CONFLICT (hotel_id, email) 
DO UPDATE SET
    full_name = EXCLUDED.full_name,
    phone = EXCLUDED.phone,
    email = EXCLUDED.email
RETURNING id;
