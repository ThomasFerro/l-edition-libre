package events

type ManuscriptReviewed struct{}

func (event ManuscriptReviewed) String() string {
	return "ManuscriptReviewed{}"
}

func (event ManuscriptReviewed) ManuscriptEventName() string {
	return "ManuscriptReviewed"
}
