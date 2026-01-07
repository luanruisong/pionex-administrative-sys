package utils

import (
	"os"
)

func Env(key, defVal string) string {
	if t, ok := os.LookupEnv(key); ok {
		return t
	}
	return defVal
}
