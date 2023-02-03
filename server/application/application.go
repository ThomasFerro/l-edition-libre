package application

import (
	"fmt"

	"github.com/ThomasFerro/l-edition-libre/commands"
	"github.com/ThomasFerro/l-edition-libre/events"
	"github.com/ThomasFerro/l-edition-libre/queries"
	"golang.org/x/exp/slog"
)

// TODO: Interfacer ?
type Application struct {
	history []events.Event
}

// TODO: repasser Application en immutable après avoir persisté les évènements ailleurs
func (app *Application) manageCommandReturn(newEvents []events.Event, err error) ([]events.Event, error) {
	if err != nil {
		return nil, err
	}
	app.history = append(app.history, newEvents...)
	return newEvents, nil
}

func (app *Application) Send(command commands.Command) ([]events.Event, error) {
	slog.Info("receiving command", "command", command)
	switch typedCommand := command.(type) {
	case commands.SubmitManuscript:
		return app.manageCommandReturn(commands.HandleSubmitManuscript(typedCommand))
	case commands.CancelManuscriptSubmission:
		return app.manageCommandReturn(commands.HandleCancelManuscriptSubmission(typedCommand))
	default:
		return nil, fmt.Errorf("unmanaged command type %T", command)
	}
}

// TODO: Remplacer le retour par du générique ?
func (app Application) Query(query queries.Query) (interface{}, error) {
	switch typedQuery := query.(type) {
	case queries.ManuscriptStatus:
		// TODO: N'envoyer que l'historique pour le manuscrit demandé
		return queries.GetManuscriptStatus(app.history, typedQuery)
	default:
		return nil, fmt.Errorf("unmanaged query type %T", query)
	}
}

func NewApplication() Application {
	return Application{}
}
