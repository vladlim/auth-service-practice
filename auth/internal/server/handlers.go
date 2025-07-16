package server

import (
	"context"
	"encoding/json"
	"net/http"
)

func (s *Server) pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func (s *Server) getPeopleHandler(w http.ResponseWriter, r *http.Request) {
	people, err := s.provider.GetPeople(context.Background())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong: " + err.Error()))
		return
	}
	body, err := json.Marshal(people)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Something went wrong"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func (s *Server) getPersonHandler(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) addPersonHandler(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) updatePersonHandler(w http.ResponseWriter, r *http.Request) {
}

func (s *Server) deletePersonHandler(w http.ResponseWriter, r *http.Request) {
}
