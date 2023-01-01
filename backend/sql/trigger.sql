CREATE OR REPLACE FUNCTION before_story_insert() RETURNS trigger AS $before_story_insert$
  DECLARE
    page_id int;
  BEGIN
    INSERT INTO page(action, author, body)
    VALUES('Welcome', NEW.author, 'Here it all begins...')
    RETURNING id INTO page_id;

    NEW.first_page_id := page_id;
    RETURN NEW;
  END;
$before_story_insert$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER before_story_insert
BEFORE INSERT ON story
FOR EACH ROW
EXECUTE FUNCTION before_story_insert();

---------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION after_story_delete() RETURNS trigger AS $after_story_delete$
  BEGIN
    DELETE FROM page
    WHERE id = OLD.first_page_id;
    RETURN OLD;
  END;
$after_story_delete$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER after_story_delete
AFTER DELETE ON story
FOR EACH ROW
EXECUTE FUNCTION after_story_delete();

---------------------------------------------------------------------------------------------

CREATE OR REPLACE FUNCTION after_choices_delete() RETURNS trigger AS $after_choices_delete$
  BEGIN
    DELETE FROM page
    WHERE id = OLD.path_id;
    RETURN OLD;
  END;
$after_choices_delete$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER after_choices_delete
AFTER DELETE ON choices
FOR EACH ROW
EXECUTE FUNCTION after_choices_delete();
