package helpers

import (
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/domain"
)

type HttpErrorMessage struct {
	Error string `json:"error"`
}

func ManageErrorAsJson(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	errorMessage := ExtractErrorMessage(err)
	WriteJson(w, errorMessage)
}

func ExtractErrorMessage(err error) HttpErrorMessage {
	typedDomainError, isDomainError := err.(domain.DomainError)
	errorMessage := HttpErrorMessage{
		Error: err.Error(),
	}
	if isDomainError {
		errorMessage = HttpErrorMessage{
			Error: typedDomainError.Name(),
		}
	}
	return errorMessage
}
