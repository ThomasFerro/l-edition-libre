package dtos

import (
	"encoding/json"
	"fmt"

	"github.com/ThomasFerro/l-edition-libre/events"
)

type ManuscriptSubmittedDto struct {
	Title    string
	Author   string
	FileName string
	FileURL  string
}

func ToDto(manuscriptEvent events.ManuscriptEvent) (interface{}, error) {
	switch manuscriptEvent.ManuscriptEventName() {
	case "ManuscriptReviewed":
		return nil, nil
	case "ManuscriptSubmissionCanceled":
		return nil, nil
	case "ManuscriptSubmitted":
		manuscriptSubmitted := manuscriptEvent.(events.ManuscriptSubmitted)
		return ManuscriptSubmittedDto{
			Title:    manuscriptSubmitted.Title,
			Author:   manuscriptSubmitted.Author,
			FileName: manuscriptSubmitted.FileName,
			FileURL:  manuscriptSubmitted.FileURL.String(),
		}, nil
	}
	return nil, fmt.Errorf("unmanaged manuscript event %v", manuscriptEvent.ManuscriptEventName())
}

func ToPayload(manuscriptEvent events.ManuscriptEvent) ([]byte, error) {
	dto, err := ToDto(manuscriptEvent)
	if err != nil {
		return nil, err
	}
	return json.Marshal(dto)
}
