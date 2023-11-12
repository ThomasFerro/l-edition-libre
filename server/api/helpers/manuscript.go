package helpers

import (
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/contexts"
)

func GetManuscriptID(r *http.Request) contexts.ManuscriptID {
	return contexts.MustParseManuscriptID(FromUrlParams(r.Context(), ":manuscriptID"))
}
