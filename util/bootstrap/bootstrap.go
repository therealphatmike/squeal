package bootstrap

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/therealphatmike/squeal/util/databases"
)

func BootstrapSqueal() error {
	if _, err := databases.InitDatabasesFile(); err != nil {
		return err
	}

	userHome, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	debugLocation := userHome + "/.squeal/debug"
	debugFile := debugLocation + "/debug.log"
	if _, err := os.Stat(debugLocation); os.IsNotExist(err) {
		if err := os.Mkdir(debugLocation, os.ModePerm); err != nil {
			return err
		}

		if _, err := os.Stat(debugFile); os.IsNotExist(err) {
			if err := os.WriteFile(debugFile, []byte(""), 0644); err != nil {
				return err
			}
		}

		f, err := tea.LogToFile(debugFile, "debug")
		if err != nil {
			return err
		}

		if err := f.Close(); err != nil {
			return err
		}
	}

	return nil
}
