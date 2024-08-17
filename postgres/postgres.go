package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/therealphatmike/squeal/util/databases"
)

func OpenDb(database databases.Database) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		database.Host,
		database.Port,
		database.Username,
		database.Password,
		database.DefaultDatabase,
	)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	defer db.Close()
	return db, nil
}
