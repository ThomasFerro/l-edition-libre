package helpers

import (
	"encoding/json"
	"net/http"
)

func WriteJson(w http.ResponseWriter, dto interface{}) {
	w.Header().Add("Content-Type", "application/json")
	manuscriptJSON, err := json.Marshal(dto)
	if err != nil {
		ManageError(w, err)
		return
	}
	w.Write(manuscriptJSON)
}
