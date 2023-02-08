package application

import (
	"fmt"

	"github.com/ThomasFerro/l-edition-libre/commands"
	"github.com/ThomasFerro/l-edition-libre/events"
	"github.com/ThomasFerro/l-edition-libre/queries"
	"golang.org/x/exp/slog"
)

type History interface {
	For(UserID, ManuscriptID) ([]events.Event, error)
	Append(UserID, ManuscriptID, []events.Event) error
}

type Application struct {
	history History
}

func (app Application) manageCommandReturn(userID UserID, manuscriptID ManuscriptID, newEvents []events.Event, err error) ([]events.Event, error) {
	if err != nil {
		return nil, err
	}
	return newEvents, app.history.Append(userID, manuscriptID, newEvents)
}

func (app Application) Send(userID UserID, manuscriptID ManuscriptID, command commands.Command) ([]events.Event, error) {
	slog.Info("receiving command", "type", fmt.Sprintf("%T", command), "manuscript_id", manuscriptID, "command", command)
	// TODO: Simplifier en jouant avec l'interface de la commande ?
	switch typedCommand := command.(type) {
	case commands.SubmitManuscript:
		newEvents, commandError := commands.HandleSubmitManuscript(typedCommand)
		return app.manageCommandReturn(userID, manuscriptID, newEvents, commandError)
	case commands.CancelManuscriptSubmission:
		history, err := app.history.For(userID, manuscriptID)
		if err != nil {
			return nil, err
		}
		newEvents, err := commands.HandleCancelManuscriptSubmission(history, typedCommand)
		return app.manageCommandReturn(userID, manuscriptID, newEvents, err)
	default:
		return nil, fmt.Errorf("unmanaged command type %T", command)
	}
}

// TODO: mieux typer le retour (générique?)
func (app Application) Query(userID UserID, manuscriptID ManuscriptID, query queries.Query) (interface{}, error) {
	slog.Info("receiving query", "type", fmt.Sprintf("%T", query), "manuscript_id", manuscriptID, "command", query)
	history, err := app.history.For(userID, manuscriptID)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch history before managing query %T", query)
	}
	switch typedQuery := query.(type) {
	case queries.ManuscriptStatus:
		return queries.GetManuscriptStatus(history, typedQuery)
	default:
		return nil, fmt.Errorf("unmanaged query type %T", query)
	}
}

func NewApplication(history History) Application {
	return Application{
		history,
	}
}
