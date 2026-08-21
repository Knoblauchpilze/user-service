
-- better-test-user@mail-client.org
INSERT INTO user_service_schema.api_user ("id", "email", "password")
  VALUES (
    '00b265e6-6638-4b1b-aeac-5898c7307eb8',
    'better-test-user@mail-client.org',
    'weakpassword'
  );

INSERT INTO user_service_schema.api_key ("id", "key", "api_user", "valid_until")
  VALUES (
    '42698272-5b8f-42db-a43c-8108eaad66e1',
    'e9c3ce0d-d6d6-45cb-ad93-c407d429469f',
    '00b265e6-6638-4b1b-aeac-5898c7307eb8',
     current_timestamp + make_interval(hours => 6)
  );
