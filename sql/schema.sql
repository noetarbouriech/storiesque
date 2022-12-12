CREATE TABLE page
(
  id BIGSERIAL PRIMARY KEY,
  title VARCHAR(32) NOT NULL,
  text VARCHAR(4096) NOT NULL
);

CREATE TABLE story
(
  id BIGSERIAL PRIMARY KEY,
  title VARCHAR(32) NOT NULL,
  description VARCHAR(512),
  first_page_id int REFERENCES page(id)
);
