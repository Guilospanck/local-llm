CREATE TABLE IF NOT EXISTS property(
   id serial PRIMARY KEY,
   color VARCHAR(50) NOT NULL,
   price NUMERIC(12, 3) NOT NULL,
   size_sqm NUMERIC(7, 2) NOT NULL
);
