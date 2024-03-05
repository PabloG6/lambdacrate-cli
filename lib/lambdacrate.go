package lib

import (
	"log"
	"os"
	"path/filepath"
)

const DefaultBaseApiURL = "https://api.lambdacrate.com"
const DefaultBaseDashboardURL = "https://lambdacrate.com"
const configFilePath = ".lambdacrate"

func DefaultConfigFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("can't find home directory")
	}
	return filepath.Join(homeDir, configFilePath)
}
