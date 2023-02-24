package persistency

import (
	"io"
	"strings"

	"github.com/ThomasFerro/l-edition-libre/commands"
	"github.com/google/uuid"
)

type inMemoryFilesSaver struct {
	Files map[string]string
}

func (saver inMemoryFilesSaver) Save(fileReader io.Reader, fileName string) (string, error) {
	fileId := uuid.New()
	builder := new(strings.Builder)
	_, err := io.Copy(builder, fileReader)
	if err != nil {
		return "", err
	}
	saver.Files[fileId.String()] = builder.String()

	return fileId.String(), nil
}

func NewFilesSaver() commands.FilesSaver {
	return inMemoryFilesSaver{
		Files: map[string]string{},
	}
}
