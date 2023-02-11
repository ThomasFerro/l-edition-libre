package helpers

import (
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/application"
)

func GetManuscriptID(r *http.Request) application.ManuscriptID {
	return application.MustParseManuscriptID(FromUrlParams(r.Context(), ":manuscriptID"))
}
