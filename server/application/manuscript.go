package application

import "github.com/ThomasFerro/l-edition-libre/domain"

// TODO: UUID / GUID ?
type ManuscriptID string

type Manuscripts interface {
	Persists(ManuscriptID, domain.Manuscript) error
}
