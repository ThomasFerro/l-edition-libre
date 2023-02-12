package commands

import (
	"fmt"

	"github.com/ThomasFerro/l-edition-libre/events"
)

type SubmitManuscript struct {
	Title  string
	Author string
}

func HandleSubmitManuscript(command SubmitManuscript) ([]events.Event, CommandError) {
	return []events.Event{
		events.ManuscriptSubmitted{
			Title:  command.Title,
			Author: command.Author,
		},
	}, nil
}

func (c SubmitManuscript) String() string {
	return fmt.Sprintf("SubmitManuscript{Title %v, Author %v}", c.Title, c.Author)
}
