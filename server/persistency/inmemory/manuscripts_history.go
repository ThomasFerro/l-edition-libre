package inmemory

import (
	"context"

	"github.com/ThomasFerro/l-edition-libre/application"
	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/utils"
)

type ManuscriptsHistory struct {
	history utils.OrderedMap[application.ManuscriptID, []application.ContextualizedEvent]
}

func (manuscripts ManuscriptsHistory) For(manuscriptID application.ManuscriptID) ([]application.ContextualizedEvent, error) {
	return manuscripts.history.Of(manuscriptID), nil
}

func (manuscripts ManuscriptsHistory) ForAll() (utils.OrderedMap[application.ManuscriptID, []application.ContextualizedEvent], error) {
	return manuscripts.history, nil
}

func (manuscripts ManuscriptsHistory) ForAllOfUser(userID application.UserID) (utils.OrderedMap[application.ManuscriptID, []application.ContextualizedEvent], error) {
	userHistory := utils.NewOrderedMap[application.ManuscriptID, []application.ContextualizedEvent]()
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
	manuscriptID := ctx.Value(contexts.ManuscriptIDContextKey{}).(application.ManuscriptID)
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
		history: utils.NewOrderedMap[application.ManuscriptID, []application.ContextualizedEvent](),
	}
}
