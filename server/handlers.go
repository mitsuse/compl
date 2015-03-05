package server

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (s *Server) getCandidates(writer http.ResponseWriter, requst *http.Request) {
	header := writer.Header()
	header.Set("Content-Type", "application/json")

	var state *State
	encodedState := requst.FormValue("state")
	decoder := json.NewDecoder(strings.NewReader(encodedState))

	if err := decoder.Decode(&state); err != nil {
		// TODO: Return error response.
		return
	}

	candSeq := s.Predictor().Predict(state.Context, state.Prefix, state.K)

	encoder := json.NewEncoder(writer)
	encoder.Encode(candSeq)
}

func (s *Server) getDescription(writer http.ResponseWriter, requst *http.Request) {
	header := writer.Header()
	header.Set("Content-Type", "application/json")

	descMap := make(map[string]interface{})
	descMap["order"] = s.Predictor().Order()

	encoder := json.NewEncoder(writer)
	encoder.Encode(descMap)
}

type State struct {
	Context []string `json:"context"`
	Prefix  string   `json:"prefix"`
	K       int      `json:"k"`
}
