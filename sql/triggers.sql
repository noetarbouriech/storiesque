CREATE OR REPLACE FUNCTION before_story_insert() RETURNS trigger AS $before_story_insert$
  DECLARE
    page_id int;
  BEGIN
    INSERT INTO page(title, text)
    VALUES('Welcome', 'Placeholder')
    RETURNING id INTO page_id;

    NEW.first_page_id := page_id;
    RETURN NEW;
  END;
$before_story_insert$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER before_story_insert
BEFORE INSERT ON story
FOR EACH ROW
EXECUTE FUNCTION before_story_insert();
