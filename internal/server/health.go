package server

import "net/http"

// health is a simple healthcheck which returns OK and 200.
func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("OK"))
}
