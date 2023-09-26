package commands

import (
	"context"
	"fmt"

	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/ThomasFerro/l-edition-libre/events"
)

type CancelManuscriptSubmission struct{}

func HandleCancelManuscriptSubmission(ctx context.Context, command Command) ([]events.Event, domain.DomainError) {
	manuscript := rehydrateFromContext(ctx)
	return manuscript.Cancel()
}

func (c CancelManuscriptSubmission) String() string {
	return fmt.Sprintf("CancelManuscriptSubmission{}")
}
