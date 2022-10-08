CREATE TABLE planets (
  id int NOT NULL AUTO_INCREMENT,
  created_at timestamp NOT NULL,
  updated_at timestamp NOT NULL,
  deleted_at timestamp,
  name varchar(100) NOT NULL,
  climates JSON NOT NULL,
  terrains JSON NOT NULL,
  PRIMARY KEY (id),
  CONSTRAINT UC_PLANET_NAME UNIQUE (name)
);