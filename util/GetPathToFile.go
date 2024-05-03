package util

import (
	"fmt"
	"os"
	"strings"
)

func GetPathToFile(localDir, filename string) (string, error) {
	pathToWorkingDir, err := os.Getwd()

	if err != nil {
		fmt.Printf("Could not get working directory")
		return "", err
	}

	projectRootDir := "investment-manager"
	pathTokens := strings.Split(pathToWorkingDir, "/")
	var pathToFile string

	for _, token := range pathTokens {
		if token == "" {
			continue
		}

		pathToFile += "/" + token

		if token == projectRootDir {
			break
		}
	}

	pathToFile += localDir + "/" + filename

	return pathToFile, nil
}
