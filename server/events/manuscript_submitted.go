package events

import "fmt"

type ManuscriptSubmitted struct {
	Title    string
	Author   string
	FileName string
	FilePath string
}

func (event ManuscriptSubmitted) String() string {
	return fmt.Sprintf("ManuscriptSubmitted{Title %v, Author %v}", event.Title, event.Author)
}

func (event ManuscriptSubmitted) ManuscriptEventName() string {
	return "ManuscriptSubmitted"
}
