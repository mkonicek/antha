package jfile

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"
	"time"
)

func TestDefaultClient(t *testing.T) {
	if len(os.Getenv("ANTHA_PASSWORD")) == 0 {
		t.Skip("missing required environment variables")
	}

	const (
		golden   = "testing"
		filename = "test_file"
	)

	c, err := DefaultClient()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	files, err := c.ListFiles(ctx, "")
	if err != nil {
		t.Error(err)
	}

	seen := make(map[string]bool)
	for _, f := range files {
		seen[f.Name] = true
	}

	for _, n := range []string{"input.json", "output.json"} {
		if !seen[n] {
			t.Errorf("expecting to find file %s but it was not found", n)
		}
	}

	w := c.NewWriter(ctx, filename)
	defer w.Close() // nolint

	if _, err := io.Copy(w, bytes.NewReader([]byte(golden))); err != nil {
		t.Error(err)
	}

	if err := w.Close(); err != nil {
		t.Error(err)
	}

	// Spin waiting for file to appear
	found := false
	for i := 0; i < 5 && !found; i++ {
		<-time.After(1 * time.Second)
		files, err := c.ListFiles(ctx, "")
		if err != nil {
			t.Fatal(err)
		}
		for _, f := range files {
			if f.Name == filename {
				found = true
				break
			}
		}
	}

	if !found {
		t.Fatalf("file not found")
	}

	r := c.NewReader(ctx, "", filename)
	defer r.Close() // nolint

	var out bytes.Buffer
	if _, err := io.Copy(&out, r); err != nil {
		t.Error(err)
	}

	if e, f := golden, out.String(); e != f {
		t.Errorf("expecting %q found %q", e, f)
	}
}
