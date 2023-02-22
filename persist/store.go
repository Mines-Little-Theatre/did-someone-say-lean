package persist

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Mines-Little-Theatre/did-someone-say-lean/persist/queries"
	_ "modernc.org/sqlite"
)

type Store interface {
	// CheckBypassRateLimit checks whether a message sent by the user in the channel should bypass rate limiting
	CheckBypassRateLimit(userId, channelId string) (bool, error)
	// PollRateLimit checks whether enough time has passed since the last rate-limited event for either ID, and if so, registers a new rate-limited event as a side effect.
	PollRateLimit(userId, channelId string) (bool, error)
	// Close should be called before the application exits
	Close() error
}

const expectedSchemaVersion = 1

type dbStore struct {
	*sql.DB
}

func Connect() (Store, error) {
	conn, err := sql.Open("sqlite", os.Getenv("LEAN_DB"))
	if err != nil {
		return nil, err
	}

	row := conn.QueryRow(queries.Get("get_schema_version"))
	var schemaVersion int
	err = row.Scan(&schemaVersion)
	if err != nil {
		conn.Close()
		return nil, err
	}

	if schemaVersion < expectedSchemaVersion {
		conn.Close()
		return nil, fmt.Errorf("schema version %d is too low (expected %d) -- please upgrade the database", schemaVersion, expectedSchemaVersion)
	} else if schemaVersion > expectedSchemaVersion {
		conn.Close()
		return nil, fmt.Errorf("schema version %d is too high (expected %d) -- please upgrade the application", schemaVersion, expectedSchemaVersion)
	}

	return &dbStore{DB: conn}, nil
}

func (s *dbStore) CheckBypassRateLimit(userId, channelId string) (bool, error) {
	row := s.QueryRow(queries.Get("check_bypass_rate_limit"), userId, channelId)
	var result int
	if err := row.Scan(&result); err != nil {
		return false, err
	}
	return result > 0, nil
}

func (s *dbStore) PollRateLimit(userId, channelId string) (bool, error) {
	tx, err := s.Begin()
	if err != nil {
		return false, err
	}

	row := tx.QueryRow(queries.Get("check_rate_limit"), userId, channelId)
	var result int
	if err = row.Scan(&result); err != nil {
		tx.Rollback()
		return false, err
	}

	rateLimited := result > 0
	if !rateLimited {
		_, err = tx.Exec(queries.Get("update_rate_limit"), userId, channelId)
		if err != nil {
			tx.Rollback()
			return false, err
		}
	}

	return result > 0, tx.Commit()
}
