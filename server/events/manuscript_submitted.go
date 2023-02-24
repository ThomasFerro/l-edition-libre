package events

import (
	"fmt"
	"net/url"
)

type ManuscriptSubmitted struct {
	Title    string
	Author   string
	FileName string
	FileURL  url.URL
}

func (event ManuscriptSubmitted) String() string {
	return fmt.Sprintf("ManuscriptSubmitted{Title %v, Author %v, FileName %v, FileURL %v}", event.Title, event.Author, event.FileName, event.FileURL)
}

func (event ManuscriptSubmitted) ManuscriptEventName() string {
	return "ManuscriptSubmitted"
}
