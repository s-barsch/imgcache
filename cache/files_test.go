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
