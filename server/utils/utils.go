package utils

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

func Check(err error) {
	if err != nil {
		slog.Error(fmt.Sprint(err))
		panic(err)
	}
}

func CurrDir() string {
	execPath, err := os.Executable()
	if err != nil {
		slog.Error(fmt.Sprint(err))
		return ""
	}
	return filepath.Dir(execPath)
}
