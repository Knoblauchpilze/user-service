
-- another-test-user@another-provider.com
INSERT INTO user_service_schema.api_user ("id", "email", "password")
  VALUES (
    '4f26321f-d0ea-46a3-83dd-6aa1c6053aaf',
    'another-test-user@another-provider.com',
    'super-strong-password'
  );

INSERT INTO user_service_schema.api_key ("id", "key", "api_user", "valid_until")
  VALUES (
    'fd8136c4-c584-4bbf-a390-53d5c2548fb8',
    '2da3e9ec-7299-473a-be0f-d722d870f51a',
    '4f26321f-d0ea-46a3-83dd-6aa1c6053aaf',
     current_timestamp + make_interval(hours => 6)
  );
