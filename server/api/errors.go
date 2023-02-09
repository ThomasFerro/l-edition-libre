package api

import (
	"fmt"
	"net/http"

	"github.com/ThomasFerro/l-edition-libre/commands"
)

type HttpErrorMessage struct {
	Error string `json:"error"`
}

func manageError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	typedCommandError, isCommandError := err.(commands.CommandError)
	errorMessage := HttpErrorMessage{
		Error: err.Error(),
	}
	fmt.Printf("\n\nisCommandError?%v\n", isCommandError)
	if isCommandError {
		errorMessage = HttpErrorMessage{
			Error: typedCommandError.Name(),
		}
	}
	writeJson(w, errorMessage)
}
