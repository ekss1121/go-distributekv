package web

import (
	"distributeKV/db"
	"fmt"
	"net/http"
)

type Server struct {
	db *db.Database
}

func CreateServer(db *db.Database) *Server {
	return &Server{
		db: db,
	}
}

func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	value, _ := s.db.Get(key)
	fmt.Fprintf(w, "%s: %s", key, value)
}

func (s *Server) SetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	value := r.Form.Get("value")
	s.db.Set(key, []byte(value))
	updated, _ := s.db.Get(key)
	fmt.Fprintf(w, "%s: %s", key, updated)
}

func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, nil)
}
