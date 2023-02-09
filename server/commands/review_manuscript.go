package commands

import (
	"fmt"

	"github.com/ThomasFerro/l-edition-libre/events"
)

type ReviewManuscript struct {
}

func HandleReviewManuscript(command ReviewManuscript) ([]events.Event, CommandError) {
	// TODO: Vérifier le status
	return []events.Event{
		events.ManuscriptReviewed{},
	}, nil
}

func (c ReviewManuscript) String() string {
	return fmt.Sprintf("ReviewManuscript{}")
}
