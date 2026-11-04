package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func (s *Server) campaignVideo(w http.ResponseWriter, r *http.Request) {
	value, err := s.campaigns.Get(r.Context(), r.PathValue("id"))
	number, parseErr := strconv.Atoi(r.PathValue("part"))
	if err != nil || parseErr != nil || number < 1 || number > len(value.Parts) {
		http.NotFound(w, r)

		return
	}

	part := value.Parts[number-1]
	path := part.File
	if id := r.URL.Query().Get("take"); id != "" {
		path = ""

		for _, take := range part.Takes {
			if take.ID == id {
				path = take.File
			}
		}
	}

	base, err := filepath.Abs(filepath.Join(s.config.Workspace, "exports"))
	absolute, pathErr := filepath.Abs(path)
	if err != nil || pathErr != nil || path == "" || filepath.Ext(path) != ".mp4" {
		http.NotFound(w, r)

		return
	}

	relative, err := filepath.Rel(base, absolute)
	if err != nil {
		http.NotFound(w, r)

		return
	}

	root, err := os.OpenRoot(base)
	if err != nil {
		http.NotFound(w, r)

		return
	}
	defer root.Close()

	file, err := root.Open(relative)
	if err != nil {
		http.NotFound(w, r)

		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil || !stat.Mode().IsRegular() {
		http.NotFound(w, r)

		return
	}

	w.Header().Set("Content-Type", "video/mp4")
	http.ServeContent(w, r, filepath.Base(path), stat.ModTime(), file)
}
