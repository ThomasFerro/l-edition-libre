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

func ToEvents(toMap []ContextualizedEvent) []events.Event {
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

// TODO: A terme, il faudrait que l'app ne s'occupe que de faire du dispatch de commandes / queries et laisse
// les history à la couche au-dessus
type Application struct {
	ManuscriptsHistory ManuscriptsHistory
	UsersHistory       UsersHistory
	managedCommands    ManagedCommands
	managedQueries     ManagedQueries
}

// TODO: Déplacer dans un middleware post traitement
func (app Application) manageCommandReturn(ctx context.Context, newEvents []events.Event, err error) ([]events.Event, error) {
	if err != nil {
		return nil, err
	}
	userID := ctx.Value(contexts.UserIDContextKey).(UserID)
	contextualizedUserEvents := []ContextualizedEvent{}
	contextualizedManuscriptEvents := []ContextualizedEvent{}
	for _, nextEvent := range newEvents {
		if _, isUserEvent := nextEvent.(events.UserEvent); isUserEvent {
			contextualizedUserEvents = append(contextualizedUserEvents, ContextualizedEvent{
				Event: nextEvent,
				Context: EventContext{
					UserID: userID,
				},
			})
			continue
		}
		if _, isManuscriptEvent := nextEvent.(events.ManuscriptEvent); isManuscriptEvent {
			contextualizedManuscriptEvents = append(contextualizedManuscriptEvents, ContextualizedEvent{
				Event: nextEvent,
				Context: EventContext{
					UserID: userID,
				},
			})
		}
	}
	if len(contextualizedUserEvents) != 0 {
		err = app.UsersHistory.Append(ctx, contextualizedUserEvents)
		if err != nil {
			return nil, err
		}
	}
	if len(contextualizedManuscriptEvents) != 0 {
		err = app.ManuscriptsHistory.Append(ctx, contextualizedManuscriptEvents)
		if err != nil {
			return nil, err
		}
	}
	return newEvents, nil
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
	history, err := app.ManuscriptsHistory.For(manuscriptID)
	if err != nil {
		return nil, err
	}
	switch typedQuery := query.(type) {
	case queries.ManuscriptStatus:
		return queries.GetManuscriptStatus(ToEvents(history), typedQuery)
	default:
		return nil, fmt.Errorf("unmanaged query type %T", query)
	}
}

func (app Application) ManuscriptsQuery(userID UserID, query queries.Query) (interface{}, error) {
	slog.Info("receiving query", "type", fmt.Sprintf("%T", query), "query", query)
	switch typedQuery := query.(type) {
	case queries.ManuscriptsToReview:
		history, err := app.ManuscriptsHistory.ForAll()
		if err != nil {
			return nil, err
		}
		return queries.GetManuscriptsToReview(toEventsByManuscript(history), typedQuery)
	case queries.WriterManuscripts:
		history, err := app.ManuscriptsHistory.ForAllOfUser(userID)
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
