package db

import (
	"github.com/KnoblauchPilze/user-service/pkg/db/postgresql"
)

var dbTestConfig = postgresql.NewConfigForLocalhost("test_db", "test_user", "test_password")
