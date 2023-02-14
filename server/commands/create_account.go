package commands

import (
	"context"
	"fmt"

	"github.com/ThomasFerro/l-edition-libre/events"
)

type CreateAccount struct {
	DisplayedName string
}

func HandleCreateAccount(ctx context.Context, command Command) ([]events.Event, CommandError) {
	return []events.Event{
		events.AccountCreated{},
	}, nil
}

func (c CreateAccount) String() string {
	return fmt.Sprintf("CreateAccount{DisplayedName %v}", c.DisplayedName)
}
