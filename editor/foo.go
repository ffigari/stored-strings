package editor

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"os/exec"
)

type ReadFile struct {
	bytes []byte
}

func (f ReadFile) Bytes() []byte {
	return f.bytes
}

var ErrReachedMaxFileSize = errors.New("reached max file size of 3MB")

func NewFile(src multipart.File) (*ReadFile, error) {
	defer src.Close()

	bytes, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	return &ReadFile{
		bytes: bytes,
	}, nil
}

type Image struct {
	file ReadFile
}

func NewImage(file *ReadFile) (*Image, error) {
	if file == nil {
		return nil, errors.New("nil read file")
	}

	return &Image{
		file: *file,
	}, nil
}

func (i *Image) Bytes() []byte {
	return i.file.Bytes()
}

type SmallImage struct {
	image Image
}

func downscaleImage(image Image, maximumSize int) (*SmallImage, error) {
	return nil, errors.New(`TODO: Implement downsizing, maybe using https://github.com/disintegration/imaging`)
}

func NewSmallImage(image Image) (*SmallImage, error) {
	return downscaleImage(image, 3 * 1024 * 1024)
}

func FindMatchingFeatures(image []SmallImage) {
	return
}

func main() {
	input1 := []byte(`foo
	asd
	qwe`)
	input2 := []byte(`asdqwe`)
	output, err := call([][]byte{
		input1, input2,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(output)
}

func call(input [][]byte) (string, error) {
	cmd := exec.Command("./venv/bin/python", "foo.py")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", err
	}

	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return "", err
	}

	for _, data := range input {
		if err := binary.Write(stdin, binary.BigEndian, uint32(len(data))); err != nil {
			return "", err
		}

		if _, err := stdin.Write(data); err != nil {
			return "", err
		}
	}
	
	stdin.Close()

	if err := cmd.Wait(); err != nil {
		return "", err
	}

	return stdout.String(), nil
}
