package config

import (
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func PrependHomeDirectory(path string) string {
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(currentUser.HomeDir, path)
}

func PrependConfigDirectory(path string) string {
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		return filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "roadie", path)
	} else {
		currentUser, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}

		return filepath.Join(currentUser.HomeDir, ".config/roadie", path)
	}
}

func ReadPasswordFile(path string) (string, error) {
	passwordBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", nil
	}

	return strings.TrimSpace(string(passwordBytes)), nil
}
