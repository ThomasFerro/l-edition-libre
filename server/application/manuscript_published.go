package application

import (
	"context"

	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/events"
)

func HandleManuscriptReviewed(ctx context.Context, manuscriptReviewed events.Event) (context.Context, []events.Event, error) {
	manuscriptID := contexts.FromContext[ManuscriptID](ctx, contexts.ManuscriptIDContextKey{})
	ctx = context.WithValue(ctx, contexts.PublicationIDContextKey{}, PublicationID(manuscriptID))
	return ctx, []events.Event{
		events.PublicationMadeAvailable{},
	}, nil
}
