package helpers

import (
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/domain"
)

type HttpErrorMessage struct {
	Error string `json:"error"`
}

func ManageError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	typedDomainError, isDomainError := err.(domain.DomainError)
	errorMessage := HttpErrorMessage{
		Error: err.Error(),
	}
	if isDomainError {
		errorMessage = HttpErrorMessage{
			Error: typedDomainError.Name(),
		}
	}
	WriteJson(w, errorMessage)
}
