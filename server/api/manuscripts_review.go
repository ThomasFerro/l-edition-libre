package api

import (
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api/helpers"
	"github.com/ThomasFerro/l-edition-libre/api/html"
	"github.com/ThomasFerro/l-edition-libre/api/middlewares"
	"github.com/ThomasFerro/l-edition-libre/commands"
	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/ThomasFerro/l-edition-libre/queries"
	"golang.org/x/exp/slog"
)

func handleManuscriptReviewSubmission(w http.ResponseWriter, r *http.Request) *http.Request {
	app := middlewares.ApplicationFromRequest(r)
	ctx, err := app.SendCommand(r.Context(), commands.ReviewManuscript{})
	if err != nil {
		slog.Error("manuscript submission review request error", err)
		helpers.ManageError(w, err)
		return r
	}
	r = r.WithContext(ctx)
	slog.Info("manuscript submission reviewed")
	helpers.WriteJson(w, "")
	return r
}

type ManuscriptToReviewDto struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

type ManuscriptsToReviewDto struct {
	Manuscripts []ManuscriptToReviewDto `json:"manuscripts"`
}

func fromDomain(manuscripts []domain.Manuscript) ManuscriptsToReviewDto {
	dto := ManuscriptsToReviewDto{
		Manuscripts: make([]ManuscriptToReviewDto, 0),
	}

	for _, manuscript := range manuscripts {
		dto.Manuscripts = append(dto.Manuscripts, ManuscriptToReviewDto{
			Title:  manuscript.Title,
			Author: manuscript.Author,
		})
	}

	return dto
}

func handleGetManuscriptsToReview(w http.ResponseWriter, r *http.Request) *http.Request {
	app := middlewares.ApplicationFromRequest(r)
	queryResult, err := app.Query(r.Context(), queries.ManuscriptsToReview{})
	if err != nil {
		slog.Error("manuscripts to review request error", err)
		helpers.ManageError(w, err)
		return r
	}
	manuscripts, castedSuccessfuly := queryResult.([]domain.Manuscript)
	if !castedSuccessfuly {
		slog.Error("manuscripts to review query result casting error", err)
		helpers.ManageError(w, err)
		return r
	}

	return html.RespondWithTemplate(w, r, fromDomain(manuscripts), "manuscripts-review.gohtml")
}
