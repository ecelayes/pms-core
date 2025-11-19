TRUNCATE inventory, room_types, rate_plans, hotels CASCADE;

WITH new_hotel AS (
    INSERT INTO hotels (name, subdomain, currency)
    VALUES ('Hotel Demo Go', 'demo-go', 'USD')
    RETURNING id
),
new_plans AS (
    INSERT INTO rate_plans (hotel_id, name, active)
    SELECT id, 'Tarifa Estándar', true FROM new_hotel
    RETURNING id, hotel_id
),
new_rooms AS (
    INSERT INTO room_types (hotel_id, name, code, base_price)
    SELECT id, 'Habitación Estándar', 'STD', 100.00 FROM new_hotel
    RETURNING id, hotel_id, code 
),
new_suite AS (
    INSERT INTO room_types (hotel_id, name, code, base_price)
    SELECT id, 'Suite Presidencial', 'SUI', 250.00 FROM new_hotel
    RETURNING id, hotel_id
)
INSERT INTO inventory (
    hotel_id, room_type_id, rate_plan_id, date, 
    total_inventory, booked_count, price
)
SELECT 
    r.hotel_id,
    r.id as room_type_id,
    p.id as rate_plan_id,
    CURRENT_DATE + (i || ' days')::interval as date,
    10 as total_inventory,
    0 as booked_count,
    CASE 
        WHEN r.code = 'SUI' THEN 250.00
        ELSE 100.00 
    END as price
FROM generate_series(0, 30) i
CROSS JOIN new_rooms r
CROSS JOIN new_plans p
UNION ALL
SELECT 
    s.hotel_id,
    s.id as room_type_id,
    p.id as rate_plan_id,
    CURRENT_DATE + (i || ' days')::interval as date,
    2 as total_inventory,
    0 as booked_count,
    250.00 as price
FROM generate_series(0, 30) i
CROSS JOIN new_suite s
CROSS JOIN new_plans p;
