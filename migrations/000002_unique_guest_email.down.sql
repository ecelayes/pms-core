ALTER TABLE guests 
DROP CONSTRAINT IF EXISTS guests_hotel_id_email_key;

INSERT INTO guests 
SELECT * FROM _backup_guests_v2
ON CONFLICT (id) DO NOTHING;

DROP TABLE IF EXISTS _backup_guests_v2;
