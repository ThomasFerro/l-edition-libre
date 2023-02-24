package contexts

import (
	"context"
)

type UserIDContextKey struct{}
type NewEventsContextKey struct{}
type ManuscriptIDContextKey struct{}
type PublicationIDContextKey struct{}
type ApplicationContextKey struct{}
type UsersHistoryContextKey struct{}
type ContextualizedUserHistoryContextKey struct{}
type ManuscriptsHistoryContextKey struct{}
type ContextualizedManuscriptsHistoryContextKey struct{}
type ContextualizedManuscriptHistoryContextKey struct{}
type PublicationsHistoryContextKey struct{}
type ContextualizedPublicationHistoryContextKey struct{}
type FilesSaverContextKey struct{}

func FromContext[T any](ctx context.Context, key interface{}) T {
	return ctx.Value(key).(T)
}

func FromContextOrDefault[T any](ctx context.Context, key interface{}, defaultValue T) T {
	value := ctx.Value(key)
	if value == nil {
		return defaultValue
	}
	return value.(T)
}
