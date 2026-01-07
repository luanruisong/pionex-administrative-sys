package utils

import (
	"fmt"
	"os"
)

func DirExists(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("目录不存在: %s", path)
		}
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("路径存在但不是目录: %s", path)
	}
	return nil
}

func Mkdir(path string) error {
	return os.MkdirAll(path, 0755)
}

func TryMkdir(path string) error {
	if err := DirExists(path); err != nil {
		return Mkdir(path)
	}
	return nil
}
