package commands

import (
	"context"
	"fmt"

	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/ThomasFerro/l-edition-libre/events"
)

type CancelManuscriptSubmission struct{}

func HandleCancelManuscriptSubmission(ctx context.Context, command Command) ([]events.Event, CommandError) {
	history := contexts.ManuscriptHistoryFromContext(ctx)
	manuscript := domain.Rehydrate(history)
	if manuscript.Status != domain.PendingReview {
		return nil, AManuscriptShouldBePendingReviewForItsSubmissionToBeCanceled{
			actualStatus: manuscript.Status,
		}
	}
	return []events.Event{
		events.ManuscriptSubmissionCanceled{},
	}, nil
}

func (c CancelManuscriptSubmission) String() string {
	return fmt.Sprintf("CancelManuscriptSubmission{}")
}

type AManuscriptShouldBePendingReviewForItsSubmissionToBeCanceled struct {
	actualStatus domain.Status
}

func (commandError AManuscriptShouldBePendingReviewForItsSubmissionToBeCanceled) Error() string {
	return fmt.Sprintf("manuscript should be pending review for its subscription to be canceled (%v)", commandError.actualStatus)
}

func (commandError AManuscriptShouldBePendingReviewForItsSubmissionToBeCanceled) Name() string {
	return "AManuscriptShouldBePendingReviewForItsSubmissionToBeCanceled"
}
