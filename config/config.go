package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	ApiKeyPath    string
	isInitialized bool
}

var config = Config{
	ApiKeyPath:    "",
	isInitialized: false,
}

func GetConfig() (*Config, error) {
	if !config.isInitialized {
		err := initConfig()
		if err != nil {
			return nil, err
		}
		return &config, nil
	}
	return &config, nil
}

func initConfig() error {
	pathToWorkingDir, err := os.Getwd()

	if err != nil {
		fmt.Printf("Could not get working directory")
		return err
	}

	projectRootDir := "investment-manager"
	pathTokens := strings.Split(pathToWorkingDir, "/")
	var pathToEnvFile string

	for _, token := range pathTokens {
		if token == "" {
			continue
		}

		pathToEnvFile += "/" + token

		if token == projectRootDir {
			break
		}
	}

	pathToEnvFile += "/.env"

	err = godotenv.Load(pathToEnvFile)

	if err != nil {
		fmt.Printf("Error loading .env\n%v\n", err)
		return err
	}

	config.ApiKeyPath = os.Getenv("API_KEY_PATH")

	if config.ApiKeyPath == "" {
		return errors.New("API Key Path Not Set")
	}

	config.isInitialized = true
	return nil
}
