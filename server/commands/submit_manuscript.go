package commands

import (
	"context"
	"fmt"

	"github.com/ThomasFerro/l-edition-libre/events"
)

type SubmitManuscript struct {
	Title  string
	Author string
}

func HandleSubmitManuscript(ctx context.Context, command Command) ([]events.Event, CommandError) {
	submitManuscript := command.(SubmitManuscript)
	return []events.Event{
		events.ManuscriptSubmitted{
			Title:  submitManuscript.Title,
			Author: submitManuscript.Author,
		},
	}, nil
}

func (c SubmitManuscript) String() string {
	return fmt.Sprintf("SubmitManuscript{Title %v, Author %v}", c.Title, c.Author)
}
