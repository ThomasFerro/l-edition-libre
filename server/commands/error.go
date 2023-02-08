package commands

import "fmt"

type CommandError interface {
	Name() string
	Error() string
}

type ManuscriptNotFound struct{}

func (commandError ManuscriptNotFound) Error() string {
	return fmt.Sprintf("resource not found")
}

func (commandError ManuscriptNotFound) Name() string {
	return "ManuscriptNotFound"
}
