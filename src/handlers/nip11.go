package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/feanor306/nostressgo/src/logger"
)

type Nip11obj struct {
	Name          string
	Description   string
	Pubkey        string   // admin pubkey
	Contact       string   // alternative contact like mailto or https
	SupportedNips []string `json:"supported_nips"`
	Software      string   // git repo
	Version       string
}

func getResponseObj() ([]byte, error) {
	no := Nip11obj{
		Name:          "NOSTRessGO",
		Description:   "NOSTR relay written in golang.",
		SupportedNips: []string{"1", "9", "11", "12", "14", "16", "20", "33", "40", "50"},
		Software:      "https://github.com/Feanor306/nostressgo",
		Version:       "1.0.0",
	}

	return json.MarshalIndent(no, "", "  ")
}

func Nip11response(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/nostr+json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")

	response, err := getResponseObj()

	if err != nil {
		logger.New().Error().Err(err).Msg("unable to respond to nip11")
		fmt.Fprint(w, string(err.Error()))
		return
	}

	fmt.Fprint(w, string(response))
}
