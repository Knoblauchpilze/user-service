
CREATE TABLE user_group (
  id UUID NOT NULL,
  api_user UUID NOT NULL,
  name TEXT NOT NULL,
  description TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  FOREIGN KEY (api_user) REFERENCES api_user(id)
);

CREATE INDEX user_group_api_user_index ON user_group (api_user);
