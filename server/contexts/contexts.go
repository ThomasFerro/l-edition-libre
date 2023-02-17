package contexts

import (
	"context"
)

const UserIDContextKey = "UserID"
const NewEventsContextKey = "NewEvents"
const ManuscriptIDContextKey = "ManuscriptID"
const ApplicationContextKey = "Application"
const UsersHistoryContextKey = "UserHistory"
const ContextualizedUserHistoryContextKey = "ContextualizedUserHistory"
const ManuscriptsHistoryContextKey = "ManuscriptsHistory"
const ContextualizedManuscriptsHistoryContextKey = "ContextualizedManuscriptsHistory"
const ContextualizedManuscriptHistoryContextKey = "ContextualizedManuscriptHistory"

func FromContext[T any](ctx context.Context, key string) T {
	return ctx.Value(key).(T)
}

func FromContextOrDefault[T any](ctx context.Context, key string, defaultValue T) T {
	value := ctx.Value(key)
	if value == nil {
		return defaultValue
	}
	return value.(T)
}
