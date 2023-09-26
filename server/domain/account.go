package domain

import (
	"fmt"

	"github.com/ThomasFerro/l-edition-libre/events"
)

type Account struct {
	DisplayedName AccountDisplayedName
}

type AccountDisplayedName string

func (a Account) String() string {
	return fmt.Sprintf("Account{DisplayedName %v}", a.DisplayedName)
}

func PromoteToEditor() ([]events.Event, DomainError) {
	return []events.Event{
		events.UserPromotedToEditor{},
	}, nil
}

func CreateAccount(displayedName string) ([]events.Event, DomainError) {
	return []events.Event{
		events.AccountCreated{},
	}, nil
}
