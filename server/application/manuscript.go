package application

import (
	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/ThomasFerro/l-edition-libre/events"
)

type Manuscripts interface {
	Persists(events.ManuscriptID, domain.Manuscript) error
}
