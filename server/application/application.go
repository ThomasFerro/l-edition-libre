package application

import (
	"context"
	"fmt"

	"github.com/ThomasFerro/l-edition-libre/commands"
	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/events"
	"github.com/ThomasFerro/l-edition-libre/queries"
	"golang.org/x/exp/slog"
)

// TODO: un history générique ?
type EventContext struct {
	UserID
}
type ContextualizedEvent struct {
	Event        events.Event
	Context      EventContext
	ManuscriptID ManuscriptID
}

func toEvents(toMap []ContextualizedEvent) []events.Event {
	returned := make([]events.Event, 0)
	for _, nextEvent := range toMap {
		returned = append(returned, nextEvent.Event)
	}
	return returned
}

func toEventsByManuscript(toMap map[ManuscriptID][]ContextualizedEvent) [][]events.Event {
	returned := make([][]events.Event, 0)
	for _, nextManuscript := range toMap {
		mappedEvents := make([]events.Event, 0)
		for _, nextEvent := range nextManuscript {
			mappedEvents = append(mappedEvents, nextEvent.Event)
		}
		returned = append(returned, mappedEvents)
	}
	return returned
}

type UsersHistory interface {
	For(UserID) ([]ContextualizedEvent, error)
	Append(UserID, []ContextualizedEvent) error
}
type ManuscriptsHistory interface {
	For(ManuscriptID) ([]ContextualizedEvent, error)
	ForAll() (map[ManuscriptID][]ContextualizedEvent, error)
	ForAllOfUser(UserID) (map[ManuscriptID][]ContextualizedEvent, error)
	Append(ManuscriptID, []ContextualizedEvent) error
}

type Application struct {
	manuscriptsHistory ManuscriptsHistory
	usersHistory       UsersHistory
}

func (app Application) manageManuscriptCommandReturn(ctx context.Context, manuscriptID ManuscriptID, newEvents []events.Event, err error) ([]events.Event, error) {
	if err != nil {
		return nil, err
	}

	contextualizedEvents := make([]ContextualizedEvent, len(newEvents))
	for index, nextEvent := range newEvents {
		contextualizedEvents[index] = ContextualizedEvent{
			Event: nextEvent,
			Context: EventContext{
				UserID: ctx.Value(contexts.UserIDContextKey).(UserID),
			},
		}
	}
	return newEvents, app.manuscriptsHistory.Append(manuscriptID, contextualizedEvents)
}

func (app Application) manageUserCommandReturn(userID UserID, newEvents []events.Event, err error) ([]events.Event, error) {
	if err != nil {
		return nil, err
	}
	contextualizedEvents := make([]ContextualizedEvent, len(newEvents))
	for index, nextEvent := range newEvents {
		contextualizedEvents[index] = ContextualizedEvent{
			Event: nextEvent,
		}
	}
	return newEvents, app.usersHistory.Append(userID, contextualizedEvents)
}

func (app Application) UserHaveAccessToManuscript(userID UserID, manuscriptID ManuscriptID) (bool, error) {
	isAnEditor, err := userIsAnEditor(app.usersHistory, userID)
	if err != nil {
		slog.Warn("user role check error", "user_id", userID, "manuscript_id", manuscriptID, "error", err)
		return false, err
	}
	if isAnEditor {
		return true, nil
	}
	isManuscriptWriter, err := isTheManuscriptWriter(app.manuscriptsHistory, userID, manuscriptID)
	if err != nil {
		slog.Warn("user's link to manuscript check error", "user_id", userID, "manuscript_id", manuscriptID, "error", err)
		return false, err
	}
	if isManuscriptWriter {
		return true, nil
	}
	slog.Warn("user is not the writer nor an editor", "user_id", userID, "manuscript_id", manuscriptID)
	return false, nil
}

func (app Application) UserIsAnEditor(userID UserID) (bool, error) {
	return userIsAnEditor(app.usersHistory, userID)
}

// TODO: Simplifier en jouant avec l'interface de la commande ? Sans avoir à envoyer d'id
func (app Application) SendUserCommand(userID UserID, command commands.Command) ([]events.Event, error) {
	slog.Info("receiving command", "type", fmt.Sprintf("%T", command), "user_id", userID, "command", command)
	switch typedCommand := command.(type) {
	case commands.CreateAccount:
		newEvents, commandError := commands.HandleCreateAccount(typedCommand)
		return app.manageUserCommandReturn(userID, newEvents, commandError)
	case commands.PromoteUserToEditor:
		newEvents, commandError := commands.HandlePromoteUserToEditor(typedCommand)
		return app.manageUserCommandReturn(userID, newEvents, commandError)
	default:
		return nil, fmt.Errorf("unmanaged command type %T", command)
	}
}

func (app Application) SendManuscriptCommand(ctx context.Context, manuscriptID ManuscriptID, command commands.Command) ([]events.Event, error) {
	slog.Info("receiving command", "type", fmt.Sprintf("%T", command), "manuscript_id", manuscriptID, "command", command)

	history, err := app.manuscriptsHistory.For(manuscriptID)
	if err != nil {
		return nil, err
	}
	eventsHistory := toEvents(history)

	switch typedCommand := command.(type) {
	case commands.SubmitManuscript:
		newEvents, commandError := commands.HandleSubmitManuscript(typedCommand)
		return app.manageManuscriptCommandReturn(ctx, manuscriptID, newEvents, commandError)
	case commands.ReviewManuscript:
		newEvents, commandError := commands.HandleReviewManuscript(eventsHistory, typedCommand)
		return app.manageManuscriptCommandReturn(ctx, manuscriptID, newEvents, commandError)
	case commands.CancelManuscriptSubmission:
		newEvents, err := commands.HandleCancelManuscriptSubmission(eventsHistory, typedCommand)
		return app.manageManuscriptCommandReturn(ctx, manuscriptID, newEvents, err)
	default:
		return nil, fmt.Errorf("unmanaged command type %T", command)
	}
}

// TODO: mieux typer le retour (générique?)
func (app Application) ManuscriptQuery(manuscriptID ManuscriptID, query queries.Query) (interface{}, error) {
	slog.Info("receiving query", "type", fmt.Sprintf("%T", query), "manuscript_id", manuscriptID, "query", query)
	history, err := app.manuscriptsHistory.For(manuscriptID)
	if err != nil {
		return nil, err
	}
	switch typedQuery := query.(type) {
	case queries.ManuscriptStatus:
		return queries.GetManuscriptStatus(toEvents(history), typedQuery)
	default:
		return nil, fmt.Errorf("unmanaged query type %T", query)
	}
}

func (app Application) ManuscriptsQuery(userID UserID, query queries.Query) (interface{}, error) {
	slog.Info("receiving query", "type", fmt.Sprintf("%T", query), "query", query)
	switch typedQuery := query.(type) {
	case queries.ManuscriptsToReview:
		history, err := app.manuscriptsHistory.ForAll()
		if err != nil {
			return nil, err
		}
		return queries.GetManuscriptsToReview(toEventsByManuscript(history), typedQuery)
	case queries.WriterManuscripts:
		history, err := app.manuscriptsHistory.ForAllOfUser(userID)
		if err != nil {
			return nil, err
		}
		return queries.GetWriterManuscripts(toEventsByManuscript(history), typedQuery)
	default:
		return nil, fmt.Errorf("unmanaged query type %T", query)
	}
}

func NewApplication(manuscriptsHistory ManuscriptsHistory, usersHistory UsersHistory) Application {
	return Application{
		manuscriptsHistory,
		usersHistory,
	}
}
