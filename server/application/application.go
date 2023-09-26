package application

import (
	"context"
	"fmt"

	"github.com/ThomasFerro/l-edition-libre/commands"
	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/ThomasFerro/l-edition-libre/events"
	"github.com/ThomasFerro/l-edition-libre/queries"
	"golang.org/x/exp/slog"
)

type commandType string
type commandHandler func(context.Context, commands.Command) ([]events.Event, domain.DomainError)
type ManagedCommands = map[commandType]commandHandler

type eventType string
type eventHandler func(context.Context, events.Event) (context.Context, []events.Event, domain.DomainError)
type ManagedEvents = map[eventType]eventHandler

type queryType string
type queryHandler func(context.Context, queries.Query) (interface{}, error)
type ManagedQueries = map[queryType]queryHandler

type Application struct {
	managedCommands ManagedCommands
	managedEvents   ManagedEvents
	managedQueries  ManagedQueries
}

func (app Application) applyEventsHandlers(ctx context.Context, commandEvent events.Event) (context.Context, []events.Event, error) {
	eventType := eventType(fmt.Sprintf("%T", commandEvent))
	if eventHandler, found := app.managedEvents[eventType]; found {
		return eventHandler(ctx, commandEvent)
	}
	return ctx, []events.Event{}, nil
}

func (app Application) manageCommandReturn(ctx context.Context, commandEvents []events.Event, err error) (context.Context, error) {
	newEvents := commandEvents

	for _, commandEvent := range commandEvents {
		var eventsToAppend []events.Event
		ctx, eventsToAppend, err = app.applyEventsHandlers(ctx, commandEvent)
		if err != nil {
			return ctx, err
		}
		newEvents = append(newEvents, eventsToAppend...)
	}

	return context.WithValue(ctx, contexts.NewEventsContextKey{}, newEvents), err
}

func (app Application) SendCommand(ctx context.Context, command commands.Command) (context.Context, error) {
	sentCommandType := commandType(fmt.Sprintf("%T", command))
	slog.Info("receiving command", "type", string(sentCommandType))

	if commandHandler, exists := app.managedCommands[sentCommandType]; exists {
		newEvents, err := commandHandler(ctx, command)
		return app.manageCommandReturn(ctx, newEvents, err)
	}
	return nil, fmt.Errorf("unhandled command %v", sentCommandType)
}

// TODO: Queries = consommation de projections ?
func (app Application) Query(ctx context.Context, query queries.Query) (interface{}, error) {
	sentQueryType := queryType(fmt.Sprintf("%T", query))
	slog.Info("receiving query", "type", string(sentQueryType))

	if queryHandler, exists := app.managedQueries[sentQueryType]; exists {
		return queryHandler(ctx, query)
	}
	return nil, fmt.Errorf("unhandled query %v", sentQueryType)
}

func NewApplication(managedCommands ManagedCommands, managedEvents ManagedEvents, managedQueries ManagedQueries) Application {
	return Application{
		managedCommands,
		managedEvents,
		managedQueries,
	}
}
