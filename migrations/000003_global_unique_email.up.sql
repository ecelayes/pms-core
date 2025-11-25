ALTER TABLE users DROP CONSTRAINT IF EXISTS users_hotel_id_email_key;

ALTER TABLE users ADD CONSTRAINT users_email_key UNIQUE (email);
