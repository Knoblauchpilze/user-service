
-- user1
INSERT INTO user_service_schema.api_user ("id", "email", "password")
  VALUES (
    '0463ed3d-bfc9-4c10-b6ee-c223bbca0fab',
    'user1@mail.com',
    'pwd1'
  );

-- https://www.postgresql.org/docs/current/functions-datetime.html#FUNCTIONS-DATETIME-CURRENT
INSERT INTO user_service_schema.api_key ("id", "key", "api_user", "valid_until")
  VALUES (
    'a5eff7a9-9bd6-4f51-9b42-a7ca5ffd3f5e',
    '3e8d49a3-9220-4ea0-88eb-299520c6ab85',
    '0463ed3d-bfc9-4c10-b6ee-c223bbca0fab',
     current_timestamp + make_interval(hours => 6)
  );
