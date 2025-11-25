-- name: CreateHotel :one
INSERT INTO hotels (name, subdomain, currency, owner_id)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: ListHotelsByOwner :many
SELECT * FROM hotels WHERE owner_id = $1;
