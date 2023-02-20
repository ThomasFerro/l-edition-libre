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

type commandType string
type commandHandler func(context.Context, commands.Command) ([]events.Event, commands.CommandError)
type ManagedCommands = map[commandType]commandHandler
type queryType string
type queryHandler func(context.Context, queries.Query) (interface{}, error)
type ManagedQueries = map[queryType]queryHandler

type Application struct {
	managedCommands ManagedCommands
	managedQueries  ManagedQueries
}

func (app Application) manageCommandReturn(ctx context.Context, newEvents []events.Event, err error) (context.Context, error) {
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

// TODO: Queries = consommation de projections
func (app Application) Query(ctx context.Context, query queries.Query) (interface{}, error) {
	sentQueryType := queryType(fmt.Sprintf("%T", query))
	slog.Info("receiving query", "type", string(sentQueryType))

	if queryHandler, exists := app.managedQueries[sentQueryType]; exists {
		return queryHandler(ctx, query)
	}
	return nil, fmt.Errorf("unhandled query %v", sentQueryType)
}

func NewApplication(managedCommands ManagedCommands, managedQueries ManagedQueries) Application {
	return Application{
		managedCommands,
		managedQueries,
	}
}
