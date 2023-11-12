package application

import (
	"context"

	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/events"
	"golang.org/x/exp/slog"
)

func IsAnEditor(ctx context.Context) (bool, error) {
	getHistory := contexts.FromContextOrDefault(ctx, contexts.ContextualizedUserHistoryContextKey{}, func(c context.Context) ([]ContextualizedEvent, error) {
		return []ContextualizedEvent{}, nil
	})
	history, err := getHistory(ctx)

	if err != nil {
		return false, err
	}
	for _, nextEvent := range history {
		_, isAUserEditorEvent := nextEvent.Event().(events.UserPromotedToEditor)
		if isAUserEditorEvent {
			return true, nil
		}
	}

	return false, nil
}

func UserHaveAccessToManuscript(ctx context.Context) (bool, error) {
	isAnEditor, err := IsAnEditor(ctx)
	if err != nil {
		slog.Warn("user role check error", "error", err)
		return false, err
	}
	if isAnEditor {
		return true, nil
	}
	isManuscriptWriter, err := isTheManuscriptWriter(ctx)
	if err != nil {
		slog.Warn("user's link to manuscript check error", "error", err)
		return false, err
	}
	if isManuscriptWriter {
		return true, nil
	}
	slog.Warn("user is not the writer nor an editor")
	return false, nil
}

type UsersHistory interface {
	For(contexts.UserID) ([]ContextualizedEvent, error)
	Append(context.Context, []ContextualizedEvent) error
}
