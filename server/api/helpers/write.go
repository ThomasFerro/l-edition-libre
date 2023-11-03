package helpers

import (
	"encoding/json"
	"net/http"
)

// TODO: Supprimer quand tout sera en HTML
func WriteJson(w http.ResponseWriter, dto interface{}) {
	w.Header().Add("Content-Type", "application/json")
	manuscriptJSON, err := json.Marshal(dto)
	if err != nil {
		ManageErrorAsJson(w, err)
		return
	}
	w.Write(manuscriptJSON)
}
