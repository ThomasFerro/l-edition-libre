package application

import (
	"fmt"

	"github.com/ThomasFerro/l-edition-libre/commands"
	"github.com/ThomasFerro/l-edition-libre/events"
	"github.com/ThomasFerro/l-edition-libre/queries"
	"golang.org/x/exp/slog"
)

type History interface {
	For(ManuscriptID) ([]events.Event, error)
	Append(ManuscriptID, []events.Event) error
}

type Application struct {
	history History
}

func (app Application) manageCommandReturn(manuscriptID ManuscriptID, newEvents []events.Event, err error) ([]events.Event, error) {
	if err != nil {
		return nil, err
	}
	return newEvents, app.history.Append(manuscriptID, newEvents)
}

func (app Application) Send(manuscriptID ManuscriptID, command commands.Command) ([]events.Event, error) {
	slog.Info("receiving command", "type", fmt.Sprintf("%T", command), "manuscript_id", manuscriptID, "command", command)
	switch typedCommand := command.(type) {
	case commands.SubmitManuscript:
		newEvents, err := commands.HandleSubmitManuscript(typedCommand)
		return app.manageCommandReturn(manuscriptID, newEvents, err)
	case commands.CancelManuscriptSubmission:
		newEvents, err := commands.HandleCancelManuscriptSubmission(typedCommand)
		return app.manageCommandReturn(manuscriptID, newEvents, err)
	default:
		return nil, fmt.Errorf("unmanaged command type %T", command)
	}
}

// TODO: Remplacer le retour par du générique ?
func (app Application) Query(manuscriptID ManuscriptID, query queries.Query) (interface{}, error) {
	slog.Info("receiving query", "type", fmt.Sprintf("%T", query), "manuscript_id", manuscriptID, "command", query)
	history, err := app.history.For(manuscriptID)
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
