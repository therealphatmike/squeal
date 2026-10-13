package databases

import (
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type DatabaseFile struct {
	Databases []Database
}

func InitDatabasesFile() (bool, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return false, err
	}

	dbConfigDir := userHome + "/.squeal"
	dbConfigFile := dbConfigDir + "/databases.toml"

	if _, err := os.Stat(dbConfigDir); os.IsNotExist(err) {
		if err := os.Mkdir(dbConfigDir, os.ModePerm); err != nil {
			return false, err
		}
	}

	if _, err := os.Stat(dbConfigFile); os.IsNotExist(err) {
		if err := os.WriteFile(dbConfigFile, []byte(""), 0644); err != nil {
			return false, err
		}
	}

	return true, nil
}

func AddDatabaseConnection(newDb Database) error {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dbConfigDir := userHome + "/.squeal"
	dbConfigFile := dbConfigDir + "/databases.toml"

	if _, err := os.Stat(dbConfigFile); os.IsNotExist(err) {
		if err := os.WriteFile(dbConfigFile, []byte(""), 0644); err != nil {
			log.Panic("unable to create file")
			return err
		}
	}

	f, err := os.OpenFile(dbConfigFile, os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Panic("unable to read file")
		return err
	}

	dbFile := DatabaseFile{}
	if _, err := toml.DecodeFile(dbConfigFile, &dbFile); err != nil {
		log.Panic(fmt.Sprintf("Unable to unmarshal database file: %s", err))
		return err
	}

	dbs := dbFile.Databases
	dbs = append(dbs, newDb)

	newContent := DatabaseFile{
		Databases: dbs,
	}
	if err := toml.NewEncoder(f).Encode(newContent); err != nil {
		log.Panic(fmt.Sprintf("unable to write to file: %s", err))
		return err
	}
	if err := f.Close(); err != nil {
		log.Panic("unable to close file")
		return err
	}

	return nil
}

func DeleteDatabaseConnection(connectionName string) error {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	dbConfigDir := userHome + "/.squeal"
	dbConfigFile := dbConfigDir + "/databases.toml"

	if _, err := os.Stat(dbConfigFile); os.IsNotExist(err) {
		return fmt.Errorf("database config file not found at %s", dbConfigFile)
	}

	f, err := os.OpenFile(dbConfigFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("unable to open file: %w", err)
	}
	defer f.Close()

	dbFile := DatabaseFile{}
	if _, err := toml.DecodeFile(dbConfigFile, &dbFile); err != nil {
		return fmt.Errorf("unable to unmarshal database file: %w", err)
	}

	dbs := dbFile.Databases
	found := false
	var filtered []Database
	for _, db := range dbs {
		if db.ConnectionName == connectionName {
			found = true
		} else {
			filtered = append(filtered, db)
		}
	}

	if !found {
		return fmt.Errorf("database connection '%s' not found", connectionName)
	}

	newContent := DatabaseFile{
		Databases: filtered,
	}
	if err := toml.NewEncoder(f).Encode(newContent); err != nil {
		return fmt.Errorf("unable to write to file: %w", err)
	}

	return nil
}

func ReadDatabaseConfigs() ([]Database, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dbConfigDir := userHome + "/.squeal"
	dbConfigFile := dbConfigDir + "/databases.toml"

	dbFile := DatabaseFile{}
	if _, err := toml.DecodeFile(dbConfigFile, &dbFile); err != nil {
		log.Panic(fmt.Sprintf("Unable to unmarshal database file: %s", err))
		return nil, err
	}

	return dbFile.Databases, nil
}
