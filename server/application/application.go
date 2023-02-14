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

/* TODO: Chantier de simplification de la partie application
- Lancer l'application avec une map[commandType]commandHandler / map[queryType]queryHandler
- Une seule méthode Command et une seule Query
- Les commandes et queries ont toutes la même API
- Les données supplémentaires sont à chercher dans le context, qui pourra potentiellement être raffiné en fonction de la commande / query (limiter l'accès en fonction de l'utilisateur)
- Helpers aussi disponibles dans le context (ApplicationContext ?)
- Queries = consommation de projections
*/

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

type CommandType string
type CommandHandler func(context.Context, commands.Command) ([]events.Event, commands.CommandError)
type ManagedCommands = map[CommandType]CommandHandler
type queryType string
type queryHandler func() interface{}
type ManagedQueries = map[queryType]queryHandler

type Application struct {
	manuscriptsHistory ManuscriptsHistory
	usersHistory       UsersHistory
	managedCommands    ManagedCommands
	managedQueries     ManagedQueries
}

func (app Application) manageCommandReturn(ctx context.Context, newEvents []events.Event, err error) ([]events.Event, error) {
	if err != nil {
		return nil, err
	}
	userID := ctx.Value(contexts.UserIDContextKey).(UserID)
	contextualizedEvents := make([]ContextualizedEvent, len(newEvents))
	for index, nextEvent := range newEvents {
		contextualizedEvents[index] = ContextualizedEvent{
			Event: nextEvent,
			Context: EventContext{
				UserID: userID,
			},
		}
	}
	// TODO: Comment gérer où on envoie les events ?
	err = app.usersHistory.Append(ctx, contextualizedEvents)
	if err != nil {
		return nil, err
	}
	return newEvents, app.manuscriptsHistory.Append(ctx, contextualizedEvents)
}

func (app Application) SendCommand(ctx context.Context, command commands.Command) ([]events.Event, error) {
	sentCommandType := CommandType(fmt.Sprintf("%T", command))
	slog.Info("receiving command", "type", string(sentCommandType))

	if commandHandler, exists := app.managedCommands[sentCommandType]; exists {
		newEvents, err := commandHandler(ctx, command)
		return app.manageCommandReturn(ctx, newEvents, err)
	}
	return nil, fmt.Errorf("unhandled command %v", sentCommandType)
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

func NewApplication(manuscriptsHistory ManuscriptsHistory, usersHistory UsersHistory, managedCommands ManagedCommands, managedQueries ManagedQueries) Application {
	return Application{
		manuscriptsHistory,
		usersHistory,
		managedCommands,
		managedQueries,
	}
}
