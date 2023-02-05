package events

type ManuscriptSubmitted struct{}

func (event ManuscriptSubmitted) String() string {
	return "ManuscriptSubmitted{}"
}
