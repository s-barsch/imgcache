package cache

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOriginalsLeaveSeekPreviewsAlone(t *testing.T) {
	root := t.TempDir()

	if err := os.MkdirAll(filepath.Join(root, "seek"), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"cover.jpg", "seek/200831_223743.jpg"} {
		if err := os.WriteFile(filepath.Join(root, name), nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	got, err := getOriginals(root)
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 1 || got[0].base() != "cover.jpg" {
		t.Fatalf("got %v, want only cover.jpg", got)
	}
}

func TestDeleteCachedClearsSeekPreviewCaches(t *testing.T) {
	root := t.TempDir()

	files := map[string]bool{
		"cover.jpg":                           true,
		"img/1280/cover.jpg":                  true,
		"img/dims/cover.jpg.txt":              true,
		"seek/200831_223743.jpg":              true,
		"seek/img/1280/200831_223743.jpg":     false,
		"seek/img/1280/200831_223743.avif":    false,
		"seek/img/dims/200831_223743.jpg.txt": false,
	}
	for name := range files {
		p := filepath.Join(root, name)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, nil, 0o644); err != nil {
			t.Fatal(err)
		}
	}

	if err := DeleteCached(root); err != nil {
		t.Fatal(err)
	}
	if err := DeleteEmptyFolders(root); err != nil {
		t.Fatal(err)
	}

	for name, kept := range files {
		if exists(filepath.Join(root, name)) != kept {
			t.Errorf("%v: kept = %v, want %v", name, !kept, kept)
		}
	}
	if exists(filepath.Join(root, "seek/img")) {
		t.Error("seek/img should be gone")
	}
}
