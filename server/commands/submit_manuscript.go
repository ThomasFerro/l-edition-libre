package commands

import (
	"context"
	"fmt"
	"io"

	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/events"
	"github.com/ThomasFerro/l-edition-libre/ports"
)

type SubmitManuscript struct {
	Title    string
	Author   string
	File     io.Reader
	FileName string
}

type UnableToPersistFile struct {
	FileName   string
	InnerError error
}

func (commandError UnableToPersistFile) Error() string {
	return commandError.InnerError.Error()
}

func (commandError UnableToPersistFile) Name() string {
	return "UnableToPersistFile"
}

func HandleSubmitManuscript(ctx context.Context, command Command) ([]events.Event, CommandError) {
	submitManuscript := command.(SubmitManuscript)

	filesSaver := contexts.FromContext[ports.FilesSaver](ctx, contexts.FilesSaverContextKey{})

	fileURL, err := filesSaver.Save(submitManuscript.File, submitManuscript.FileName)
	if err != nil {
		return nil, UnableToPersistFile{
			FileName:   submitManuscript.FileName,
			InnerError: err,
		}
	}
	return []events.Event{
		events.ManuscriptSubmitted{
			Title:    submitManuscript.Title,
			Author:   submitManuscript.Author,
			FileName: submitManuscript.FileName,
			FileURL:  fileURL,
		},
	}, nil
}

func (c SubmitManuscript) String() string {
	return fmt.Sprintf("SubmitManuscript{Title %v, Author %v, FileName %v}", c.Title, c.Author, c.FileName)
}
