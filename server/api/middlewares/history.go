package middlewares

import (
	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/events"
	"github.com/ThomasFerro/l-edition-libre/utils"
)

func mapHistories[T comparable](original utils.OrderedMap[T, []application.ContextualizedEvent]) map[T][]events.DecoratedEvent {
	returned := map[T][]events.DecoratedEvent{}

	for _, keyValuePair := range original.Map() {
		returned[keyValuePair.Key] = mapHistory(keyValuePair.Value)
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
