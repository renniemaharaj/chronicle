package imageassets

import (
	"bytes"
	"image"
	"image/png"
	"io"
	"testing"
)

func TestImportedImageTransportRejectsNonImagesAndPathTraversal(t *testing.T) {
	root := t.TempDir()
	if _, err := Put(root, []byte("not an image")); err == nil {
		t.Fatal("non-image upload accepted")
	}

	var content bytes.Buffer
	if err := png.Encode(&content, image.NewRGBA(image.Rect(0, 0, 3, 5))); err != nil {
		t.Fatal(err)
	}

	asset, err := Put(root, content.Bytes())
	if err != nil || asset.Width != 3 || asset.Height != 5 {
		t.Fatalf("image import failed: %+v %v", asset, err)
	}

	file, err := Open(root, asset.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	saved, err := io.ReadAll(file)
	if err != nil || !bytes.Equal(saved, content.Bytes()) {
		t.Fatal("imported image bytes were not preserved")
	}

	for _, path := range []string{"/etc/passwd", "/api/media/images/../../private.png", "/api/media/images/file.svg"} {
		if file, err := Open(root, path); err == nil {
			file.Close()
			t.Fatal("unconfined image path accepted")
		}
	}
}
