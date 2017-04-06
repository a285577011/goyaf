package lib

import (
	"os"
	"os/exec"
	"path/filepath"
)

func GetCurrentPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	return path
}
