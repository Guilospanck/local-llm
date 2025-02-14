CREATE TABLE IF NOT EXISTS property_views(
   id serial PRIMARY KEY,
   property_id INTEGER references property(id),
   view_id INTEGER references view(id)
);
