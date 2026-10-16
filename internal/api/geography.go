package api

import "net/http"

func (s *Server) searchPlaces(w http.ResponseWriter, r *http.Request) {
	results, err := s.geography.Search(r.Context(), r.URL.Query().Get("q"))
	if err != nil {
		fail(w, err, http.StatusBadGateway)
		return
	}

	write(w, http.StatusOK, results)
}
