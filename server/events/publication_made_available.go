package events

type PublicationMadeAvailable struct{}

func (event PublicationMadeAvailable) String() string {
	return "PublicationMadeAvailable{}"
}

func (event PublicationMadeAvailable) PublicationEventName() string {
	return "PublicationMadeAvailable"
}
