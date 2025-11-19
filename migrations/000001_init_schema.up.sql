DROP TABLE IF EXISTS booking_extras CASCADE;
DROP TABLE IF EXISTS booking_rooms CASCADE;
DROP TABLE IF EXISTS bookings CASCADE;
DROP TABLE IF EXISTS guests CASCADE;
DROP TABLE IF EXISTS inventory CASCADE;
DROP TABLE IF EXISTS extras CASCADE;
DROP TABLE IF EXISTS rate_plans CASCADE;
DROP TABLE IF EXISTS rooms CASCADE;
DROP TABLE IF EXISTS room_types CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS hotels CASCADE;

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE hotels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    subdomain VARCHAR(50) UNIQUE NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    timezone VARCHAR(50) DEFAULT 'UTC',
    settings JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hotel_id UUID REFERENCES hotels(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'receptionist',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(hotel_id, email)
);

CREATE TABLE room_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hotel_id UUID REFERENCES hotels(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(20),
    max_occupancy INT NOT NULL DEFAULT 2,
    base_price DECIMAL(10,2) NOT NULL,
    amenities JSONB DEFAULT '[]'::jsonb,
    description TEXT
);

CREATE TABLE rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hotel_id UUID REFERENCES hotels(id) ON DELETE CASCADE,
    room_type_id UUID REFERENCES room_types(id),
    number VARCHAR(20) NOT NULL,
    floor INT,
    status VARCHAR(20) DEFAULT 'clean',
    UNIQUE(hotel_id, number)
);

CREATE TABLE rate_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hotel_id UUID REFERENCES hotels(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    is_non_refundable BOOLEAN DEFAULT FALSE,
    meals_included VARCHAR(50) DEFAULT 'none',
    active BOOLEAN DEFAULT TRUE
);

CREATE TABLE extras (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hotel_id UUID REFERENCES hotels(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    pricing_model VARCHAR(30) NOT NULL DEFAULT 'per_person_per_night', 
    price DECIMAL(10, 2) NOT NULL,
    active BOOLEAN DEFAULT TRUE
);

CREATE TABLE inventory (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hotel_id UUID NOT NULL REFERENCES hotels(id) ON DELETE CASCADE,
    room_type_id UUID NOT NULL REFERENCES room_types(id),
    rate_plan_id UUID REFERENCES rate_plans(id),
    date DATE NOT NULL,
    total_inventory INT NOT NULL,
    booked_count INT DEFAULT 0,
    hold_count INT DEFAULT 0,
    price DECIMAL(10, 2) NOT NULL,
    is_closed BOOLEAN DEFAULT FALSE,
    min_stay INT DEFAULT 1,
    version BIGINT DEFAULT 1,
    UNIQUE(room_type_id, rate_plan_id, date)
);

CREATE INDEX idx_inventory_lookup ON inventory (hotel_id, date, room_type_id);

CREATE TABLE guests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hotel_id UUID REFERENCES hotels(id),
    email VARCHAR(255),
    full_name VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    passport_number VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hotel_id UUID REFERENCES hotels(id),
    guest_id UUID REFERENCES guests(id),
    code VARCHAR(12) UNIQUE NOT NULL,
    status VARCHAR(20) DEFAULT 'confirmed',
    room_amount DECIMAL(10, 2) NOT NULL,
    extras_amount DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
    total_amount DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'USD',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE booking_rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID REFERENCES bookings(id) ON DELETE CASCADE,
    room_type_id UUID REFERENCES room_types(id),
    rate_plan_id UUID REFERENCES rate_plans(id),
    check_in DATE NOT NULL,
    check_out DATE NOT NULL,
    price_per_night DECIMAL(10, 2) NOT NULL,
    assigned_room_id UUID REFERENCES rooms(id),
    guest_names TEXT[]
);

CREATE TABLE booking_extras (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    booking_id UUID REFERENCES bookings(id) ON DELETE CASCADE,
    extra_id UUID REFERENCES extras(id),
    "quantity" INT NOT NULL DEFAULT 1,
    unit_price DECIMAL(10, 2) NOT NULL,
    total_price DECIMAL(10, 2) NOT NULL,
    effective_date DATE
);
