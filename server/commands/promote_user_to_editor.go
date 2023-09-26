package commands

import (
	"context"
	"fmt"

	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/ThomasFerro/l-edition-libre/events"
)

type PromoteUserToEditor struct{}

func HandlePromoteUserToEditor(ctx context.Context, command Command) ([]events.Event, domain.DomainError) {
	return domain.PromoteToEditor()
}

func (c PromoteUserToEditor) String() string {
	return fmt.Sprintf("PromoteUserToEditor{}")
}
