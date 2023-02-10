package commands

import (
	"fmt"

	"github.com/ThomasFerro/l-edition-libre/events"
)

type PromoteUserToEditor struct{}

func HandlePromoteUserToEditor(command PromoteUserToEditor) ([]events.Event, CommandError) {
	return []events.Event{
		events.UserPromotedToEditor{},
	}, nil
}

func (c PromoteUserToEditor) String() string {
	return fmt.Sprintf("PromoteUserToEditor{}")
}
