package util

import "os"

func GetEnvValue(variable, def string) string {
	str := os.Getenv(variable)
	if str == "" {
		return def
	}
	return str
}
