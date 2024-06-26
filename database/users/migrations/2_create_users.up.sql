
CREATE TABLE api_user (
  id UUID NOT NULL,
  email TEXT NOT NULL,
  password TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  version INTEGER DEFAULT 0,
  PRIMARY KEY (id),
  UNIQUE (email)
);

CREATE TRIGGER trigger_api_user_updated_at
  BEFORE UPDATE OR INSERT ON api_user
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at();
