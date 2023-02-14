package commands

import (
	"context"

	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/events"
)

type Command interface{}

func historyFromContext(ctx context.Context) []events.Event {
	return ctx.Value(contexts.HistoryContextKey).([]events.Event)
}
