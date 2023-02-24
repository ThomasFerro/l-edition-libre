package commands

import (
	"context"
	"fmt"
	"io"

	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/events"
)

type SubmitManuscript struct {
	Title    string
	Author   string
	File     io.Reader
	FileName string
}

// TODO: Déplacer
type FilesSaver interface {
	Save(fileReader io.Reader, fileName string) (string, error)
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
	filesSaver := contexts.FromContext[FilesSaver](ctx, contexts.FilesSaverContextKey{})
	path, err := filesSaver.Save(submitManuscript.File, submitManuscript.FileName)
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
			FilePath: path,
		},
	}, nil
}

func (c SubmitManuscript) String() string {
	return fmt.Sprintf("SubmitManuscript{Title %v, Author %v, FileName %v}", c.Title, c.Author, c.FileName)
}
