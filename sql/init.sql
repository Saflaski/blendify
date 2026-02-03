CREATE TABLE mbids (
  id SERIAL PRIMARY KEY,
  value TEXT UNIQUE NOT NULL
);

CREATE TABLE genres (
  id SERIAL PRIMARY KEY,
  name TEXT UNIQUE NOT NULL
);

CREATE TABLE mbids_genres (
  mbid_id INT REFERENCES mbids(id) ON DELETE CASCADE,
  genre_id INT REFERENCES genres(id) ON DELETE CASCADE,
  PRIMARY KEY (mbid_id, genre_id)
);

