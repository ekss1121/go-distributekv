package web

import (
	"distributeKV/config"
	"distributeKV/db"
	"fmt"
	"net/http"
)

type Server struct {
	db        *db.Database
	partition config.Partition
	pCount    int
}

func CreateServer(db *db.Database, partition config.Partition, pCount int) *Server {
	return &Server{
		db:        db,
		partition: partition,
		pCount:    pCount,
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

func (s *Server) ListenAndServe() error {
	httpAddress := s.partition.Host
	return http.ListenAndServe(httpAddress, nil)
}
