package domain

import (
	"fmt"
	"io"
	"net/url"

	"github.com/ThomasFerro/l-edition-libre/events"
	"github.com/ThomasFerro/l-edition-libre/ports"
)

type Manuscript struct {
	Status  ManuscriptStatus
	Title   string
	Author  string
	FileURL url.URL
}

type ManuscriptStatus string

const (
	PendingReview ManuscriptStatus = "PendingReview"
	Canceled      ManuscriptStatus = "Canceled"
	Reviewed      ManuscriptStatus = "Reviewed"
)

func (m Manuscript) String() string {
	return fmt.Sprintf("Manuscript{Title %v, Author %v, Status %v}", m.Title, m.Author, m.Status)
}

func SubmitManuscript(filesSaver ports.FilesSaver, title string, author string, file io.Reader, fileName string) ([]events.Event, DomainError) {
	fileURL, err := filesSaver.Save(file, fileName)
	if err != nil {
		return nil, UnableToPersistFile{
			FileName:   fileName,
			InnerError: err,
		}
	}
	return []events.Event{
		events.ManuscriptSubmitted{
			Title:    title,
			Author:   author,
			FileName: fileName,
			FileURL:  fileURL,
		},
	}, nil
}

func (manuscript Manuscript) Cancel() ([]events.Event, DomainError) {
	if manuscript.Status != PendingReview {
		return nil, AManuscriptShouldBePendingReviewForItsSubmissionToBeCanceled{
			actualStatus: manuscript.Status,
		}
	}
	return []events.Event{
		events.ManuscriptSubmissionCanceled{},
	}, nil
}

func (manuscript Manuscript) Review() ([]events.Event, DomainError) {

	if manuscript.Status != PendingReview {
		return nil, AManuscriptShouldBePendingReviewToBeReviewed{
			actualStatus: manuscript.Status,
		}
	}

	return []events.Event{
		events.ManuscriptReviewed{},
	}, nil
}
