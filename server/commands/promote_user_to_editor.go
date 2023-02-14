package commands

import (
	"context"
	"fmt"

	"github.com/ThomasFerro/l-edition-libre/events"
)

type PromoteUserToEditor struct{}

func HandlePromoteUserToEditor(ctx context.Context, command Command) ([]events.Event, CommandError) {
	return []events.Event{
		events.UserPromotedToEditor{},
	}, nil
}

func (c PromoteUserToEditor) String() string {
	return fmt.Sprintf("PromoteUserToEditor{}")
}
