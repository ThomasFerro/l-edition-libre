package application

import (
	"fmt"

	"github.com/ThomasFerro/l-edition-libre/commands"
	"github.com/ThomasFerro/l-edition-libre/events"
	"github.com/ThomasFerro/l-edition-libre/queries"
	"golang.org/x/exp/slog"
)

// TODO: un history générique ?
type UsersHistory interface {
	For(UserID) ([]events.Event, error)
	Append(UserID, []events.Event) error
}
type ManuscriptsHistory interface {
	For(ManuscriptID) ([]events.Event, error)
	Append(ManuscriptID, []events.Event) error
}

type Application struct {
	manuscriptsHistory ManuscriptsHistory
	usersHistory       UsersHistory
}

func (app Application) manageManuscriptCommandReturn(manuscriptID ManuscriptID, newEvents []events.Event, err error) ([]events.Event, error) {
	if err != nil {
		return nil, err
	}
	return newEvents, app.manuscriptsHistory.Append(manuscriptID, newEvents)
}

func (app Application) manageUserCommandReturn(userID UserID, newEvents []events.Event, err error) ([]events.Event, error) {
	if err != nil {
		return nil, err
	}
	return newEvents, app.usersHistory.Append(userID, newEvents)
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

// TODO: Simplifier en jouant avec l'interface de la commande ? Sans avoir à envoyer d'id
func (app Application) SendUserCommand(userID UserID, command commands.Command) ([]events.Event, error) {
	slog.Info("receiving command", "type", fmt.Sprintf("%T", command), "user_id", userID, "command", command)
	switch typedCommand := command.(type) {
	case commands.CreateAccount:
		newEvents, commandError := commands.HandleCreateAccount(typedCommand)
		return app.manageUserCommandReturn(userID, newEvents, commandError)
	default:
		return nil, fmt.Errorf("unmanaged command type %T", command)
	}
}

func (app Application) SendManuscriptCommand(manuscriptID ManuscriptID, command commands.Command) ([]events.Event, error) {
	slog.Info("receiving command", "type", fmt.Sprintf("%T", command), "manuscript_id", manuscriptID, "command", command)
	switch typedCommand := command.(type) {
	case commands.SubmitManuscript:
		newEvents, commandError := commands.HandleSubmitManuscript(typedCommand)
		return app.manageManuscriptCommandReturn(manuscriptID, newEvents, commandError)
	case commands.ReviewManuscript:
		newEvents, commandError := commands.HandleReviewManuscript(typedCommand)
		return app.manageManuscriptCommandReturn(manuscriptID, newEvents, commandError)
	case commands.CancelManuscriptSubmission:
		history, err := app.manuscriptsHistory.For(manuscriptID)
		if err != nil {
			return nil, err
		}
		newEvents, err := commands.HandleCancelManuscriptSubmission(history, typedCommand)
		return app.manageManuscriptCommandReturn(manuscriptID, newEvents, err)
	default:
		return nil, fmt.Errorf("unmanaged command type %T", command)
	}
}

// TODO: mieux typer le retour (générique?)
func (app Application) Query(manuscriptID ManuscriptID, query queries.Query) (interface{}, error) {
	slog.Info("receiving query", "type", fmt.Sprintf("%T", query), "manuscript_id", manuscriptID, "command", query)
	history, err := app.manuscriptsHistory.For(manuscriptID)
	if err != nil {
		return nil, err
	}
	switch typedQuery := query.(type) {
	case queries.ManuscriptStatus:
		return queries.GetManuscriptStatus(history, typedQuery)
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
