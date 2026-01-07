package app

import (
	"os"
	"path/filepath"
	"pionex-administrative-sys/utils"
	"pionex-administrative-sys/utils/consts"
)

var (
	appHome string
	dbPath  string
	logPath string
)

func init() {
	appHome = utils.Env(consts.APP_HOME_KEY, filepath.Join(os.Getenv(consts.HOME), ".pas"))
	if err := utils.TryMkdir(appHome); err != nil {
		panic(err.Error())
	}
	dbPath = filepath.Join(appHome, "data")
	if err := utils.TryMkdir(dbPath); err != nil {
		panic(err.Error())
	}
	logPath = filepath.Join(appHome, "logs")
	if err := utils.TryMkdir(logPath); err != nil {
		panic(err.Error())
	}
}

func Home() string {
	return appHome
}

func DBPath(file string) string {
	return filepath.Join(dbPath, file)
}

func LogPath(file string) string {
	return filepath.Join(logPath, file)
}
