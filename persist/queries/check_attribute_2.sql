-- ?1: attribute name
-- ?2: an entity id
-- ?3: an entity id
-- returns one row (integer): >0 if either entity has the named attribute

SELECT EXISTS (
  SELECT * FROM attributes
  WHERE (id = $2 OR id = $3) AND attr = $1
);
