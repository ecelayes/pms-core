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
