package middlewares

import (
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/events"
)

func mapHistories[T comparable](original map[T][]application.ContextualizedEvent) [][]events.DecoratedEvent {
	returned := [][]events.DecoratedEvent{}

	for _, nextHistory := range original {
		returned = append(returned, mapHistory(nextHistory))
	}
	return returned
}

func mapHistory(original []application.ContextualizedEvent) []events.DecoratedEvent {
	mappedHistory := []events.DecoratedEvent{}
	for _, nextEvent := range original {
		mappedHistory = append(mappedHistory, nextEvent)
	}
	return mappedHistory
}
