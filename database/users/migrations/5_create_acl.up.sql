
CREATE TABLE acl (
  id UUID NOT NULL,
  user_group UUID NOT NULL,
  policy TEXT NOT NULL,
  resource TEXT NOT NULL,
  description TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  FOREIGN KEY (user_group) REFERENCES user_group(id)
);

CREATE INDEX acl_group_index ON acl (user_group);

CREATE TABLE acl_permission (
  acl UUID NOT NULL,
  permission TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (acl),
  FOREIGN KEY (acl) REFERENCES acl(id)
);
