package persist

import (
	"database/sql"
	"fmt"

	"github.com/Mines-Little-Theatre/did-someone-say-lean/persist/queries"
	"github.com/Mines-Little-Theatre/did-someone-say-lean/utils"

	_ "modernc.org/sqlite"
)

type dbStore struct {
	*sql.DB
}

const (
	applicationId uint32 = 0x4c45414e
	userVersion   uint32 = 1
)

func Connect() (Store, error) {
	db, err := sql.Open("sqlite", utils.ReadEnvRequired("LEAN_DB"))
	if err != nil {
		return nil, err
	}

	row := db.QueryRow("PRAGMA application_id;")
	var dbAppId uint32
	err = row.Scan(&dbAppId)
	if err != nil {
		db.Close()
		return nil, err
	}

	if dbAppId != applicationId && dbAppId != 0 {
		db.Close()
		return nil, fmt.Errorf("application_id mismatch: expected %d, but was %d", applicationId, dbAppId)
	}

	row = db.QueryRow("PRAGMA user_version;")
	var dbUserVer uint32
	err = row.Scan(&dbUserVer)
	if err != nil {
		db.Close()
		return nil, err
	}

	if dbUserVer > userVersion {
		db.Close()
		return nil, fmt.Errorf("user_version is too high: expected %d or lower, but was %d", userVersion, dbUserVer)
	}
	for dbUserVer < userVersion {
		_, err := db.Exec(queries.GetMigration(dbUserVer + 1))
		if err != nil {
			db.Close()
			return nil, err
		}
		dbUserVer++
	}

	return &dbStore{DB: db}, nil
}

func (s *dbStore) CheckIgnore(userId, channelId string) (bool, error) {
	return checkAttribute(s.DB, "ignore", userId, channelId)
}

func (s *dbStore) CheckBypassRateLimit(userId, channelId string) (bool, error) {
	return checkAttribute(s.DB, "bypass_rate_limit", userId, channelId)
}

func (s *dbStore) GetFallbackReaction() (string, error) {
	value, err := getProperty[string](s.DB, "fallback_reaction")
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}

func (s *dbStore) GetGigglesnortFallbackReaction() (string, error) {
	value, err := getProperty[string](s.DB, "gigglesnort_fallback_reaction")
	if err == sql.ErrNoRows {
		return "", nil
	}
	return value, err
}

func (s *dbStore) GetGigglesnortMessage(word string) (string, error) {
	row := s.QueryRow(queries.Get("get_gigglesnort_message"), word)
	var message string
	err := row.Scan(&message)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return message, err
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
