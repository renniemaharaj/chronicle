package imageassets

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"

	"chronicle/internal/domain"
)

type Image struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

func Put(workspace string, data []byte) (Image, error) {
	config, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil || (format != "png" && format != "jpeg" && format != "gif") || config.Width < 1 || config.Height < 1 || config.Width > 8192 || config.Height > 8192 || config.Width*config.Height > 40000000 {
		return Image{}, fmt.Errorf("use a PNG, JPEG or GIF image up to 8192 pixels and 40 megapixels")
	}

	extension := ".png"
	if format == "jpeg" {
		extension = ".jpg"
	} else if format == "gif" {
		extension = ".gif"
	}

	hash := sha256.Sum256(data)
	name := hex.EncodeToString(hash[:]) + extension
	directory := filepath.Join(workspace, "media", "images")
	if err := os.MkdirAll(directory, 0700); err != nil {
		return Image{}, err
	}

	root, err := os.OpenRoot(directory)
	if err != nil {
		return Image{}, err
	}
	defer root.Close()

	file, err := root.OpenFile(name, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600)
	if os.IsExist(err) {
		return Image{"/api/media/images/" + name, config.Width, config.Height}, nil
	}

	if err != nil {
		return Image{}, err
	}

	_, writeErr := file.Write(data)
	closeErr := file.Close()
	if writeErr != nil || closeErr != nil {
		_ = root.Remove(name)

		return Image{}, fmt.Errorf("image could not be saved")
	}

	return Image{"/api/media/images/" + name, config.Width, config.Height}, nil
}

func Open(workspace, address string) (*os.File, error) {
	if !domain.IsLocalImageURL(address) {
		return nil, fmt.Errorf("invalid local image reference")
	}

	root, err := os.OpenRoot(filepath.Join(workspace, "media", "images"))
	if err != nil {
		return nil, err
	}
	defer root.Close()

	return root.Open(strings.TrimPrefix(address, "/api/media/images/"))
}
