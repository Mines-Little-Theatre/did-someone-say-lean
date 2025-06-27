package persist

import (
	"database/sql"
	"fmt"

	"github.com/Mines-Little-Theatre/did-someone-say-lean/persist/queries"
)

func getProperty[T any](db *sql.DB, key string) (T, error) {
	row := db.QueryRow(queries.Get("get_property"), key)
	var value T
	err := row.Scan(&value)
	return value, err
}

// checkAttribute checks whether any identified entity has the named attribute
func checkAttribute(db *sql.DB, attr string, ids ...string) (bool, error) {
	var queryName string
	switch len(ids) {
	case 2:
		queryName = "check_attribute_2"
	default:
		return false, fmt.Errorf("no query exists to check an attribute on %d entities", len(ids))
	}

	args := make([]any, 0, len(ids)+1)
	args = append(args, attr)
	for _, id := range ids {
		args = append(args, id)
	}
	row := db.QueryRow(queries.Get(queryName), args...)
	var result bool
	err := row.Scan(&result)
	return result, err
}
