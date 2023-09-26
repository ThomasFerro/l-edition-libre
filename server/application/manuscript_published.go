package application

import (
	"context"

	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/ThomasFerro/l-edition-libre/events"
)

func HandleManuscriptReviewed(ctx context.Context, manuscriptReviewed events.Event) (context.Context, []events.Event, domain.DomainError) {
	manuscriptID := contexts.FromContext[ManuscriptID](ctx, contexts.ManuscriptIDContextKey{})
	ctx = context.WithValue(ctx, contexts.PublicationIDContextKey{}, PublicationID(manuscriptID))
	newEvents, err := domain.MakePublicationAvailable()
	return ctx, newEvents, err
}
