-- name: GetInventoryByDateRange :many
SELECT 
    id,
    hotel_id,
    room_type_id,
    date,
    total_inventory,
    booked_count,
    price,
    version
FROM inventory
WHERE 
    date >= $1 AND date <= $2
ORDER BY date ASC;

-- name: UpdateInventoryCount :one
-- This query uses Optimistic Locking (version) to avoid overbooking.
UPDATE inventory
SET 
    booked_count = booked_count + 1,
    version = version + 1
WHERE 
    id = $1 
    AND version = $2
    AND booked_count < total_inventory
RETURNING id, version, booked_count;
