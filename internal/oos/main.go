package oos

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func ReadFileAtRoot(path string) ([]byte, error) {
	// to allow using this pkg in non-root directories (eg tests for the
	// postgres server connection string)
	rootAbsolutePath, err := exec.
		Command("git", "rev-parse", "--show-toplevel").
		Output()
	if err != nil {
		return nil, err
	}

	return os.ReadFile(fmt.Sprintf(
		"%s/%s",
		strings.TrimSpace(string(rootAbsolutePath)),
		path,
	))
}
