package commands

import (
	"context"

	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/events"
)

type Command interface{}

func userHistoryFromContext(ctx context.Context) []events.Event {
	return ctx.Value(contexts.UserHistoryContextKey).([]events.Event)
}

func manuscriptHistoryFromContext(ctx context.Context) []events.Event {
	return ctx.Value(contexts.ManuscriptHistoryContextKey).([]events.Event)
}
