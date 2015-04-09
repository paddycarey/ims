package config

import (
	"os"

	"github.com/Sirupsen/logrus"
)

// GetConfigFromEnv fetches an environment variable and returns as a string.
// The second argument (a boolean), if true, will cause the function to
// fatally exit if the variable is not present in the environment.
func GetConfigFromEnv(varname string, required bool) string {
	envvar := os.Getenv(varname)
	if envvar == "" && required == true {
		logrus.WithField("key", varname).Fatal("Environment variable not found.")
	}
	return envvar
}
