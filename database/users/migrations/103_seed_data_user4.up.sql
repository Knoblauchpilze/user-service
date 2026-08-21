
-- i-dont-care-about-@security.de
INSERT INTO user_service_schema.api_user ("id", "email", "password")
  VALUES (
    'beb2a2dc-2a9f-48d6-b2ca-fd3b5ca3249f',
    'i-dont-care-about-@security.de',
    'mycatismypassword'
  );

INSERT INTO user_service_schema.api_key ("id", "key", "api_user", "valid_until")
  VALUES (
    'a610adcb-d966-4617-9f15-caf6e48b6325',
    'c64f4da4-8bc5-4e19-a038-cd8755bd07d5',
    'beb2a2dc-2a9f-48d6-b2ca-fd3b5ca3249f',
     current_timestamp + make_interval(hours => 6)
  );
