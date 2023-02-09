package events

type ManuscriptReviewed struct{}

func (event ManuscriptReviewed) String() string {
	return "ManuscriptReviewed{}"
}
