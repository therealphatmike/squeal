package postgres

import (
	"database/sql"
)

type schemaRow struct {
	schemaname string
	tablename  string
}

func GetTables(db *sql.DB) ([]Schema, error) {
	rows, err := db.Query("select schemaname, tablename from pg_catalog.pg_tables ORDER BY tablename;")
	if err != nil {
		return nil, err
	}

	schemaTableHash := map[string][]string{}
	for rows.Next() {
		schemaName := ""
		tableName := ""
		rows.Scan(&schemaName, &tableName)
		schemaTableHash[schemaName] = append(schemaTableHash[schemaName], tableName)
	}

	schemas := []Schema{}
	for schema, tables := range schemaTableHash {
		schemas = append(schemas, Schema{
			Name:   schema,
			Tables: tables,
		})
	}

	return schemas, nil
}
