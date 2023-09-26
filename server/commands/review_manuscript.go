package commands

import (
	"context"
	"fmt"

	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/ThomasFerro/l-edition-libre/events"
)

type ReviewManuscript struct{}

func HandleReviewManuscript(ctx context.Context, command Command) ([]events.Event, domain.DomainError) {
	manuscript := rehydrateFromContext(ctx)
	return manuscript.Review()
}

func (c ReviewManuscript) String() string {
	return fmt.Sprintf("ReviewManuscript{}")
}
