package inmemory

import (
	"context"

	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/utils"
)

type ManuscriptsHistory struct {
	history utils.OrderedMap[contexts.ManuscriptID, []application.ContextualizedEvent]
}

func (manuscripts ManuscriptsHistory) For(manuscriptID contexts.ManuscriptID) ([]application.ContextualizedEvent, error) {
	return manuscripts.history.Of(manuscriptID), nil
}

func (manuscripts ManuscriptsHistory) ForAll() (utils.OrderedMap[contexts.ManuscriptID, []application.ContextualizedEvent], error) {
	return manuscripts.history, nil
}

func (manuscripts ManuscriptsHistory) ForAllOfUser(userID contexts.UserID) (utils.OrderedMap[contexts.ManuscriptID, []application.ContextualizedEvent], error) {
	userHistory := utils.NewOrderedMap[contexts.ManuscriptID, []application.ContextualizedEvent]()
	for _, keyValuePair := range manuscripts.history.Map() {
		manuscriptID := keyValuePair.Key
		manuscriptHistory := keyValuePair.Value
		userManuscriptHistory := []application.ContextualizedEvent{}
		for _, manuscriptEvent := range manuscriptHistory {
			if manuscriptEvent.Context.UserID == userID {
				userManuscriptHistory = append(userManuscriptHistory, manuscriptEvent)
			}
		}
		if len(userManuscriptHistory) != 0 {
			userHistory = userHistory.Upsert(manuscriptID, userManuscriptHistory)
		}
	}
	return userHistory, nil
}

func (manuscripts ManuscriptsHistory) Append(ctx context.Context, newEvents []application.ContextualizedEvent) error {
	manuscriptID := ctx.Value(contexts.ManuscriptIDContextKey{}).(contexts.ManuscriptID)
	persistedEvents := manuscripts.history.Of(manuscriptID)
	if !manuscripts.history.HasKey(manuscriptID) {
		persistedEvents = []application.ContextualizedEvent{}
	}
	persistedEvents = append(persistedEvents, newEvents...)
	manuscripts.history = manuscripts.history.Upsert(manuscriptID, persistedEvents)
	return nil
}

func NewManuscriptsHistory() application.ManuscriptsHistory {
	return ManuscriptsHistory{
		history: utils.NewOrderedMap[contexts.ManuscriptID, []application.ContextualizedEvent](),
	}
}
