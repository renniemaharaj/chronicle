package api

import (
	"net/http"

	"chronicle/internal/campaign"
)

func (s *Server) campaignRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/campaigns/from-audio", s.audioCampaign)
	mux.HandleFunc("GET /api/campaigns/{id}/parts/{part}/video", s.campaignVideo)
	mux.HandleFunc("GET /api/campaigns", func(w http.ResponseWriter, r *http.Request) {
		result, err := s.campaigns.List(r.Context())
		if err != nil {
			fail(w, err, 500)

			return
		}

		write(w, 200, result)
	})
	mux.HandleFunc("GET /api/campaigns/{id}", func(w http.ResponseWriter, r *http.Request) {
		result, err := s.campaigns.Get(r.Context(), r.PathValue("id"))
		if err != nil {
			fail(w, err, 500)

			return
		}

		write(w, 200, result)
	})
	mux.HandleFunc("POST /api/campaigns", func(w http.ResponseWriter, r *http.Request) {
		var input campaign.CreateInput
		if err := decode(w, r, &input); err != nil {
			fail(w, err, 400)

			return
		}

		result, err := s.campaigns.Create(r.Context(), input)
		if err != nil {
			fail(w, err, 422)

			return
		}

		write(w, 201, result)
	})
	mux.HandleFunc("POST /api/campaigns/{id}/actions", func(w http.ResponseWriter, r *http.Request) {
		var input campaign.ActionInput
		if err := decode(w, r, &input); err != nil {
			fail(w, err, 400)

			return
		}

		if input.ID != r.PathValue("id") {
			fail(w, http.ErrNotSupported, 400)

			return
		}

		result, err := s.campaigns.Update(r.Context(), input)
		if err != nil {
			fail(w, err, 409)

			return
		}

		write(w, 200, result)
	})
}
