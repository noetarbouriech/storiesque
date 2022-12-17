CREATE TABLE story
(
  id BIGSERIAL PRIMARY KEY,
  title VARCHAR(32) NOT NULL,
  description VARCHAR(512),
  first_page_id BIGSERIAL REFERENCES page(id)
);
