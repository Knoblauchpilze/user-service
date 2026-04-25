#!/bin/bash

echo "Setting up users database..."
echo "You may need to enter multiple times the postgres password"

export ADMIN_PASSWORD="admin_password"
export MANAGER_PASSWORD="manager_password"
export USER_PASSWORD="user_password"

echo "Creating users..."
/bin/bash database/create_user.sh database/users

echo "Creating database..."
/bin/bash database/create_database.sh database/users

echo "Seeding values..."
export DB_PATH="database/users"
make -f database/Makefile migrate

echo "All done!"
