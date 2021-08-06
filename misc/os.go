package misc

// OS specific code goes here

import (
	"log"
	"os"
	"os/user"
	"runtime"
)

func GetConfigDir() string {

	var configDir string

	switch osName := runtime.GOOS; osName {
	case "windows":
		configDir = os.Getenv("ProgramData") + "\\"
	case "linux":
		configDir = "/etc/"
	}

	configDir += "pascollector"

	return configDir
}

func GetDatabaseDir() string {

	var dbDir string

	switch osName := runtime.GOOS; osName {
	case "windows":
		dbDir = os.Getenv("ProgramData") + "\\"
	case "linux":
		dbDir = "/var/lib/"
	}

	dbDir += "pascollector"

	return dbDir
}

func GetConfigurationFilename() string {
	return GetConfigDir() + string(os.PathSeparator) + ".pascollector.yaml"
}

func CheckUser() {
	currentUser, _ := user.Current()
	if runtime.GOOS == "linux" && currentUser.Username != "root" {
		log.Fatal("setup should be ran as root")
	}
}
