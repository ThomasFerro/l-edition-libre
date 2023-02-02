package application

import (
	"errors"

	"github.com/ThomasFerro/l-edition-libre/commands"
)

// TODO: Interfacer
type Application struct{}

func (app Application) Send(command commands.Command) error {
	return errors.New("not implemented")
}

func NewApplication() Application {
	return Application{}
}
