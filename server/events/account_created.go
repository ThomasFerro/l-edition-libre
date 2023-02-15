package events

type AccountCreated struct{}

func (event AccountCreated) String() string {
	return "AccountCreated{}"
}

func (event AccountCreated) UserEventName() string {
	return "AccountCreated"
}
