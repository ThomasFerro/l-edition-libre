package commands

type CommandError interface {
	Name() string
	Error() string
}
