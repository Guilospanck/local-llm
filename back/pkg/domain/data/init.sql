CREATE TABLE IF NOT EXISTS property(
   id serial PRIMARY KEY,
   color VARCHAR(50) NOT NULL,
   price NUMERIC(12, 3) NOT NULL,
   size_sqm NUMERIC(7, 2) NOT NULL
);

CREATE TABLE IF NOT EXISTS view(
   id serial PRIMARY KEY,
   view VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS property_views(
   id serial PRIMARY KEY,
   property_id INTEGER references property(id),
   view_id INTEGER references view(id)
);

INSERT INTO property(color, price, size_sqm)
VALUES
	('Black', 12354.23, 222),
	('Red', 77777.23, 444),
	('Marble', 888999, 33),
	('Orange', 10000, 22.23),
	('Blue', 2223.23, 77.22),
	('Yellow', 111222, 88.2),
	('Cyan', 555555555, 1233);

INSERT INTO view (view)
VALUES
	('Sea'),
	('Mountains'),
	('Lake'),
	('Nature'),
	('Forest'),
	('Buildings'),
	('River');
