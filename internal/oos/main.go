package oos

import (
	"errors"
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

type File struct {
	name string
	content []byte
}

func (f File) Content() []byte {
	return f.content
}

func (f File) Name() string {
	return f.name
}

var (
	ErrNotADir = errors.New("not a dir")
)

func ReadFiles(relativePath string) ([]File, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	absolutePath := cwd+relativePath
	info, err := os.Stat(absolutePath)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, ErrNotADir
	}

	osFiles, err := os.ReadDir(absolutePath)
	if err != nil {
		return nil, err
	}

	files := []File{}
	for _, of := range osFiles {
		fileName := of.Name()

		fileContent, err := os.ReadFile(absolutePath+"/"+fileName)
		if err != nil {
			return nil, err
		}

		files = append(files, File{
			name: fileName,
			content: fileContent,
		})
	}

	return files, nil
}
