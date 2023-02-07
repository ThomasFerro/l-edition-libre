package api

import (
	"encoding/json"
	"net/http"
)

func writeJson(w http.ResponseWriter, dto interface{}) {
	w.Header().Add("Content-Type", "application/json")
	manuscriptJSON, err := json.Marshal(dto)
	if err != nil {
		manageError(w, err)
		return
	}
	w.Write(manuscriptJSON)
}
