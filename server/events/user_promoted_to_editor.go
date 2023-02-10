package events

type UserPromotedToEditor struct{}

func (event UserPromotedToEditor) String() string {
	return "UserPromotedToEditor{}"
}
