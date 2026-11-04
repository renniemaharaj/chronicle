package api

import (
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"chronicle/internal/domain"
	"chronicle/internal/imageassets"
)

func (s *Server) storyImageRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/media/images", s.importStoryImage)
	mux.HandleFunc("GET /api/media/images/{file}", s.storyImage)
}

func (s *Server) importStoryImage(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 16<<20)
	data, err := io.ReadAll(r.Body)
	if err != nil {
		fail(w, fmt.Errorf("image exceeds the 16 MB upload limit"), 413)

		return
	}

	asset, err := imageassets.Put(s.config.Workspace, data)
	if err != nil {
		fail(w, err, 422)

		return
	}

	write(w, 201, asset)
}

func (s *Server) storyImage(w http.ResponseWriter, r *http.Request) {
	file, err := imageassets.Open(s.config.Workspace, r.URL.Path)
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

	w.Header().Set("Content-Type", "image/png")

	switch filepath.Ext(stat.Name()) {
	case ".jpg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".gif":
		w.Header().Set("Content-Type", "image/gif")
	}

	http.ServeContent(w, r, stat.Name(), stat.ModTime(), file)
}

func (s *Server) validateGeneratedImages(project domain.Project) error {
	for _, asset := range project.Media {
		if asset.Generated == nil {
			continue
		}

		file, err := imageassets.Open(s.config.Workspace, asset.URL)
		if err != nil {
			return fmt.Errorf("generated image %s must be imported into this workspace first", asset.ID)
		}

		file.Close()
	}

	return nil
}
