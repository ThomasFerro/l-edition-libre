package events

type AccountCreated struct{}

func (event AccountCreated) String() string {
	return "AccountCreated{}"
}
