package domain

import "fmt"

type DomainError interface {
	Name() string
	Error() string
}

type ManuscriptNotFound struct{}

func (domainError ManuscriptNotFound) Error() string {
	return fmt.Sprintf("resource not found")
}

func (domainError ManuscriptNotFound) Name() string {
	return "ManuscriptNotFound"
}

type AManuscriptShouldBePendingReviewForItsSubmissionToBeCanceled struct {
	actualStatus ManuscriptStatus
}

func (domainError AManuscriptShouldBePendingReviewForItsSubmissionToBeCanceled) Error() string {
	return fmt.Sprintf("manuscript should be pending review for its subscription to be canceled (%v)", domainError.actualStatus)
}

func (domainError AManuscriptShouldBePendingReviewForItsSubmissionToBeCanceled) Name() string {
	return "AManuscriptShouldBePendingReviewForItsSubmissionToBeCanceled"
}

type AManuscriptShouldBePendingReviewToBeReviewed struct {
	actualStatus ManuscriptStatus
}

func (domainError AManuscriptShouldBePendingReviewToBeReviewed) Error() string {
	return fmt.Sprintf("manuscript should be pending review to be reviewed (%v)", domainError.actualStatus)
}

func (domainError AManuscriptShouldBePendingReviewToBeReviewed) Name() string {
	return "AManuscriptShouldBePendingReviewToBeReviewed"
}

type UnableToPersistFile struct {
	FileName   string
	InnerError error
}

func (domainError UnableToPersistFile) Error() string {
	return domainError.InnerError.Error()
}

func (domainError UnableToPersistFile) Name() string {
	return "UnableToPersistFile"
}
