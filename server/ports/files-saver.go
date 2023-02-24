package ports

import (
	"io"
	"net/url"
)

type FilesSaver interface {
	Save(fileReader io.Reader, fileName string) (url.URL, error)
}
