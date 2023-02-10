package helpers

import (
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/commands"
)

type HttpErrorMessage struct {
	Error string `json:"error"`
}

func ManageError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	typedCommandError, isCommandError := err.(commands.CommandError)
	errorMessage := HttpErrorMessage{
		Error: err.Error(),
	}
	if isCommandError {
		errorMessage = HttpErrorMessage{
			Error: typedCommandError.Name(),
		}
	}
	WriteJson(w, errorMessage)
}
