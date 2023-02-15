package contexts

import (
	"context"

	"github.com/ThomasFerro/l-edition-libre/events"
)

const UserIDContextKey = "UserID"
const NewEventsContextKey = "NewEvents"
const ManuscriptIDContextKey = "ManuscriptID"
const ApplicationContextKey = "Application"
const UserHistoryContextKey = "UserHistory"
const ManuscriptsHistoryContextKey = "ManuscriptsHistory"
const ManuscriptHistoryContextKey = "ManuscriptHistory"

// TODO: Pas ouf ici non plus
func UserHistoryFromContext(ctx context.Context) []events.Event {
	return ctx.Value(UserHistoryContextKey).([]events.Event)
}

func ManuscriptHistoryFromContext(ctx context.Context) []events.Event {
	return ctx.Value(ManuscriptHistoryContextKey).([]events.Event)
}
func ManuscriptsHistoryFromContext(ctx context.Context) [][]events.Event {
	return ctx.Value(ManuscriptsHistoryContextKey).([][]events.Event)
}
