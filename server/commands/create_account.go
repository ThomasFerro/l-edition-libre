package commands

import (
	"context"
	"fmt"

	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/ThomasFerro/l-edition-libre/events"
)

type CreateAccount struct {
	DisplayedName string
}

func HandleCreateAccount(ctx context.Context, command Command) ([]events.Event, domain.DomainError) {
	return domain.CreateAccount(command.(CreateAccount).DisplayedName)
}

func (c CreateAccount) String() string {
	return fmt.Sprintf("CreateAccount{DisplayedName %v}", c.DisplayedName)
}
