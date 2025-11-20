CREATE TABLE IF NOT EXISTS _backup_guests_v2 AS 
SELECT * FROM guests;

WITH duplicates_map AS (
    SELECT 
        g.id AS victim_id,
        FIRST_VALUE(g.id) OVER (
            PARTITION BY g.hotel_id, g.email 
            ORDER BY g.created_at DESC
        ) AS master_id
    FROM guests g
)
UPDATE bookings b
SET guest_id = map.master_id
FROM duplicates_map map
WHERE b.guest_id = map.victim_id
  AND map.victim_id != map.master_id;

DELETE FROM guests
WHERE id IN (
    SELECT g.id
    FROM guests g
    JOIN (
        SELECT hotel_id, email, MAX(created_at) as last_created
        FROM guests
        GROUP BY hotel_id, email
        HAVING count(*) > 1
    ) keep ON g.hotel_id = keep.hotel_id AND g.email = keep.email
    WHERE g.created_at < keep.last_created
);

ALTER TABLE guests 
ADD CONSTRAINT guests_hotel_id_email_key UNIQUE (hotel_id, email);
