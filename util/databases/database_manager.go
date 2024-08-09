package databases

import "os"

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

func AddDatabaseConnection() (bool, error) {
	return false, nil
}
