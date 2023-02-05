package events

type ManuscriptSubmissionCanceled struct{}

func (event ManuscriptSubmissionCanceled) String() string {
	return "ManuscriptSubmissionCanceled{}"
}
