package persist

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Mines-Little-Theatre/did-someone-say-lean/persist/queries"
	"github.com/Mines-Little-Theatre/did-someone-say-lean/utils"

	_ "modernc.org/sqlite"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type dbStore struct {
	*sql.DB
}

const (
	userVersion   uint32 = 1
)

func ConnectPG(connectionString string) (*sql.DB, error){
	db, err := sql.Open("pgx", connectionString)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func ConnectSQLite(connectionString string) (*sql.DB, error){
	db, err := sql.Open("sqlite", connectionString)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func Connect() (Store, error) {
	connectionString := utils.ReadEnvRequired("LEAN_DB")

	var db *sql.DB = nil
	var err error = nil


	if strings.Contains(connectionString, "postgresql"){
		db, err = ConnectPG(connectionString)
	} else {
		// db, err = ConnectSQLite(connectionString)
        return nil, fmt.Errorf("Leanbot reqiures postgresql!")
	}


	if err != nil {
		return nil, err
	}


	row := db.QueryRow(queries.GetUserVersion())
	var dbUserVer uint32
	err = row.Scan(&dbUserVer)
	if err != nil {
		// we are going to assume that if this query fails, that we need to run all migrations.
		dbUserVer = 0
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
	var result bool
	if err = row.Scan(&result); err != nil {
		tx.Rollback()
		return false, err
	}

	if !result {
		_, err = tx.Exec(queries.Get("update_rate_limit"), userId, channelId)
		if err != nil {
			tx.Rollback()
			return false, err
		}
	}

	return result, tx.Commit()
}
