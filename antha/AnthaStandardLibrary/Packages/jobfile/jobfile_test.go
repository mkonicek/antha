package jobfile

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

// Spin waiting for file to appear
func spinUntil(ctx context.Context, c *Client, filename string, seconds int) error {
	for i := 0; i < seconds; i++ {
		<-time.After(1 * time.Second)
		files, err := c.ListFiles(ctx, "")
		if err != nil {
			return err
		}
		for _, f := range files {
			if f.Name == filename {
				return nil
			}
		}
	}

	return fmt.Errorf("file not found")
}

func TestLargeFiles(t *testing.T) {
	if len(os.Getenv("ANTHA_PASSWORD")) == 0 {
		t.Skip("missing required environment variables")
	}

	const (
		filename  = "test_large_file"
		kilobytes = 1000 * 10 // 10 MB
	)

	c, err := DefaultClient()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()

	w := c.NewWriter(ctx, filename)
	defer w.Close() // nolint

	inR, inW := io.Pipe()

	go func() {
		defer inW.Close() // nolint
		var buf [1024]byte
		for i := 0; i < kilobytes; i++ {
			if _, err := inW.Write(buf[:]); err != nil {
				inW.CloseWithError(err) // nolint
				return
			}
		}
	}()

	if _, err := io.Copy(w, inR); err != nil {
		t.Error(err)
	}

	if err := w.Close(); err != nil {
		t.Error(err)
	}

	if err := spinUntil(ctx, c, filename, 5); err != nil {
		t.Fatal(err)
	}

	r := c.NewReader(ctx, "", filename)
	defer r.Close() // nolint

	n, err := io.Copy(ioutil.Discard, r)
	if err != nil {
		t.Error(err)
	}

	if e, f := 1024*int64(kilobytes), n; e != f {
		t.Errorf("expecting %d found %d", e, f)
	}
}

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

	if err := spinUntil(ctx, c, filename, 5); err != nil {
		t.Fatal(err)
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
