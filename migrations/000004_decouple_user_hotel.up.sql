ALTER TABLE users ALTER COLUMN hotel_id DROP NOT NULL;

ALTER TABLE hotels ADD COLUMN owner_id UUID REFERENCES users(id);

UPDATE hotels h SET owner_id = u.id FROM users u WHERE u.hotel_id = h.id;
