package commands

import (
	"context"
	"fmt"

	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/ThomasFerro/l-edition-libre/events"
)

type ReviewManuscript struct {
}

func HandleReviewManuscript(ctx context.Context, command Command) ([]events.Event, CommandError) {
	history := contexts.FromContext[[]events.DecoratedEvent](ctx, contexts.ContextualizedManuscriptHistoryContextKey)

	manuscript := domain.Rehydrate(events.ToEvents(history))
	fmt.Printf("\n\n history %v  manuscript %v \n\n", history, manuscript)
	if manuscript.Status != domain.PendingReview {
		return nil, AManuscriptShouldBePendingReviewToBeReviewed{
			actualStatus: manuscript.Status,
		}
	}
	return []events.Event{
		events.ManuscriptReviewed{},
	}, nil
}

func (c ReviewManuscript) String() string {
	return fmt.Sprintf("ReviewManuscript{}")
}

type AManuscriptShouldBePendingReviewToBeReviewed struct {
	actualStatus domain.Status
}

func (commandError AManuscriptShouldBePendingReviewToBeReviewed) Error() string {
	return fmt.Sprintf("manuscript should be pending review to be reviewed (%v)", commandError.actualStatus)
}

func (commandError AManuscriptShouldBePendingReviewToBeReviewed) Name() string {
	return "AManuscriptShouldBePendingReviewToBeReviewed"
}
