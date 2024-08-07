
CREATE TABLE universe (
  id UUID NOT NULL,
  name TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  version INTEGER DEFAULT 0,
  PRIMARY KEY (id),
  UNIQUE (name)
);

CREATE TRIGGER trigger_universe_updated_at
  BEFORE UPDATE OR INSERT ON universe
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at();
