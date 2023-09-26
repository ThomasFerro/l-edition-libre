package commands

import (
	"context"
	"fmt"
	"io"

	"github.com/ThomasFerro/l-edition-libre/contexts"
	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/ThomasFerro/l-edition-libre/events"
	"github.com/ThomasFerro/l-edition-libre/ports"
)

type SubmitManuscript struct {
	Title    string
	Author   string
	File     io.Reader
	FileName string
}

func HandleSubmitManuscript(ctx context.Context, command Command) ([]events.Event, domain.DomainError) {
	submitManuscript := command.(SubmitManuscript)

	filesSaver := contexts.FromContext[ports.FilesSaver](ctx, contexts.FilesSaverContextKey{})
	return domain.SubmitManuscript(filesSaver, submitManuscript.Title, submitManuscript.Author, submitManuscript.File, submitManuscript.FileName)
}

func (c SubmitManuscript) String() string {
	return fmt.Sprintf("SubmitManuscript{Title %v, Author %v, FileName %v}", c.Title, c.Author, c.FileName)
}
