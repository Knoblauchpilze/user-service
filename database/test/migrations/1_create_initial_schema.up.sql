
SET client_encoding = 'UTF8';

SET search_path = test_db_schema, pg_catalog;
SET default_tablespace = '';

CREATE OR REPLACE FUNCTION update_updated_at() RETURNS TRIGGER AS $$
  BEGIN
    NEW.updated_at = now();
    RETURN NEW;
  END;
$$ language 'plpgsql';
