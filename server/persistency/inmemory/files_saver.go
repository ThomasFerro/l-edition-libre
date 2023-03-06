package inmemory

import (
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/ThomasFerro/l-edition-libre/ports"
	"github.com/google/uuid"
)

type localFilesSaver struct{}

func (saver localFilesSaver) Save(fileReader io.Reader, fileName string) (url.URL, error) {
	fileId := uuid.New()
	dname, err := os.MkdirTemp("", "__test_files__")
	if err != nil {
		return url.URL{}, err
	}

	newFileName := fmt.Sprintf("%v/%v%v", dname, fileId.String(), fileName)
	newFile, err := os.Create(newFileName)
	defer newFile.Close()
	if err != nil {
		return url.URL{}, err
	}
	_, err = io.Copy(newFile, fileReader)
	if err != nil {
		return url.URL{}, err
	}

	newFileURL := fmt.Sprintf("file://%v", newFileName)
	parsedURL, err := url.Parse(newFileURL)
	return *parsedURL, err
}

func NewFilesSaver() ports.FilesSaver {
	return localFilesSaver{}
}
