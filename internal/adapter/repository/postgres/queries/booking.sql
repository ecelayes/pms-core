-- name: CreateBooking :one
INSERT INTO bookings (
    hotel_id, guest_id, code, status, 
    room_amount, extras_amount, total_amount, currency
) VALUES (
    $1, $2, $3, $4, 
    $5, $6, $7, $8
)
RETURNING id, created_at;

-- name: CreateBookingRoom :exec
INSERT INTO booking_rooms (
    booking_id, room_type_id, rate_plan_id, 
    check_in, check_out, price_per_night, guest_names
) VALUES (
    $1, $2, $3, 
    $4, $5, $6, $7
);

-- name: GetBookingByID :one
SELECT 
    b.id, b.code, b.status, b.total_amount, b.currency, b.created_at,
    g.full_name as guest_name, g.email as guest_email, g.phone as guest_phone
FROM bookings b
JOIN guests g ON b.guest_id = g.id
WHERE b.id = $1 LIMIT 1;

-- name: ListBookings :many
SELECT 
    b.id, b.code, b.status, b.total_amount, b.currency, b.created_at,
    g.full_name as guest_name
FROM bookings b
JOIN guests g ON b.guest_id = g.id
WHERE 
    b.hotel_id = $1
    AND (sqlc.narg('status')::text IS NULL OR b.status = sqlc.narg('status'))
ORDER BY b.created_at DESC;

-- name: UpdateBookingStatus :exec
UPDATE bookings 
SET status = $2, updated_at = NOW()
WHERE id = $1;

-- name: GetBookingRooms :many
SELECT * FROM booking_rooms WHERE booking_id = $1;
