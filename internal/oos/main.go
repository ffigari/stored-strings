package oos

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GetRootPath() (string, error) {
	rootAbsolutePath, err := exec.
		Command("git", "rev-parse", "--show-toplevel").
		Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(rootAbsolutePath)), nil
}

func ReadFileAtRoot(path string) ([]byte, error) {
	rootPath, err := GetRootPath()
	if err != nil {
		return nil, err
	}

	return os.ReadFile(fmt.Sprintf("%s/%s", rootPath, path))
}
