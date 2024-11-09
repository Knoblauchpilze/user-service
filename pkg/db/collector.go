package db

import (
	"reflect"

	"github.com/jackc/pgx/v5"
)

func getCollectorForType[T any]() pgx.RowToFunc[T] {
	var value T

	// TODO: Test this
	if reflect.ValueOf(value).Kind() == reflect.Struct {
		return pgx.RowToStructByName[T]
	}

	return pgx.RowTo[T]
}
