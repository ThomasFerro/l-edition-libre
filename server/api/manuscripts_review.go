package api

import (
	"fmt"
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/api/helpers"
	"github.com/ThomasFerro/l-edition-libre/api/middlewares"
	"github.com/ThomasFerro/l-edition-libre/commands"
	"github.com/ThomasFerro/l-edition-libre/domain"
	"github.com/ThomasFerro/l-edition-libre/queries"
	"golang.org/x/exp/slog"
)

// type ReviewManuscriptRequestDto struct{}

func handleManuscriptReviewSubmission(w http.ResponseWriter, r *http.Request) {
	// decoder := json.NewDecoder(r.Body)
	// var dto ReviewManuscriptRequestDto
	// err := decoder.Decode(&dto)
	// slog.Info("manuscript submission review request", "user_id", userID.String(), "manuscript_id", manuscriptID.String(), "body", dto)
	// if err != nil {
	// 	slog.Error("manuscript creation request dto decoding error", err)
	// 	helpers.ManageError(w, err)
	// 	return
	// }
	manuscriptID := getManuscriptID(r)
	app := middlewares.ApplicationFromRequest(r)
	_, err := app.SendManuscriptCommand(r.Context(), manuscriptID, commands.ReviewManuscript{})
	if err != nil {
		slog.Error("manuscript submission review request error", err, "manuscript_id", manuscriptID.String())
		helpers.ManageError(w, err)
		return
	}
	slog.Info("manuscript submission reviewed", "manuscript_id", manuscriptID)
	helpers.WriteJson(w, "")
}

type ManuscriptToReviewDto struct {
	Title  string `json:"title"`
	Author string `json:"author"`
}

type ManuscriptsToReviewDto struct {
	Manuscripts []ManuscriptToReviewDto `json:"manuscripts"`
}

func fromDomain(manuscripts []domain.Manuscript) ManuscriptsToReviewDto {
	fmt.Printf("\n\n %v \n\n", manuscripts)
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

func handleGetManuscriptsToReview(w http.ResponseWriter, r *http.Request) {
	app := middlewares.ApplicationFromRequest(r)
	queryResult, err := app.ManuscriptsQuery(queries.ManuscriptsToReview{})
	if err != nil {
		slog.Error("manuscripts to review request error", err)
		helpers.ManageError(w, err)
		return
	}
	manuscripts, castedSuccessfuly := queryResult.([]domain.Manuscript)
	if !castedSuccessfuly {
		slog.Error("manuscripts to review query result casting error", err)
		helpers.ManageError(w, err)
		return
	}

	helpers.WriteJson(w, fromDomain(manuscripts))
}
