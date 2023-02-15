package events

type ManuscriptSubmissionCanceled struct{}

func (event ManuscriptSubmissionCanceled) String() string {
	return "ManuscriptSubmissionCanceled{}"
}

func (event ManuscriptSubmissionCanceled) ManuscriptEventName() string {
	return "ManuscriptSubmissionCanceled"
}
