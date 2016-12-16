package trace

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"
)

func NewTestContext() (context.Context, context.CancelFunc, DoneFunc) {
	return NewContext(context.Background())
}

func TestGoOneRead(t *testing.T) {
	ctx, cancel, allDone := NewTestContext()
	defer cancel()

	Go(ctx, func(ctx context.Context) error {
		p := Issue(ctx, "noop")
		_, err := Read(ctx, p)
		return err
	})

	select {
	case <-allDone():
		if err := ctx.Err(); err != nil {
			t.Error(err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout")
	}
}

func TestCommandSequence(t *testing.T) {
	ctx, cancel, allDone := NewTestContext()
	defer cancel()

	Go(ctx, func(ctx context.Context) error {
		for i := 0; i < 5; i += 1 {
			p := Issue(ctx, "noop")
			if _, err := Read(ctx, p); err != nil {
				return err
			}
		}
		return nil
	})

	select {
	case <-allDone():
		if err := ctx.Err(); err != nil {
			t.Error(err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout")
	}
}

func TestNestedCommandSequence(t *testing.T) {
	ctx, cancel, allDone := NewTestContext()
	defer cancel()

	Go(ctx, func(ctx context.Context) error {
		for i := 0; i < 100; i += 1 {
			pidx := i
			Go(ctx, func(ctx context.Context) error {
				for i := 0; i < 10; i += 1 {
					p := Issue(ctx, fmt.Sprintf("noop.%d.%d", pidx, i))
					if _, err := Read(ctx, p); err != nil {
						return err
					}
				}
				return nil
			})
		}
		return nil
	})

	select {
	case <-allDone():
		if err := ctx.Err(); err != nil {
			t.Error(err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout")
	}
}

func TestGoGoReadAll(t *testing.T) {
	ctx, cancel, allDone := NewTestContext()
	defer cancel()

	Go(ctx, func(ctx context.Context) error {
		for i := 0; i < 5; i += 1 {
			pidx := i
			Go(ctx, func(ctx context.Context) error {
				var promises []*Promise
				for i := 0; i < 5; i += 1 {
					p := Issue(ctx, fmt.Sprintf("noop.%d.%d", pidx, i))
					promises = append(promises, p)
				}
				if _, err := ReadAll(ctx, promises...); err != nil {
					return err
				}
				return nil
			})
		}
		return nil
	})

	select {
	case <-allDone():
		if err := ctx.Err(); err != nil {
			t.Error(err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout")
	}
}
func TestGoReadAll(t *testing.T) {
	ctx, cancel, allDone := NewTestContext()
	defer cancel()

	Go(ctx, func(ctx context.Context) error {
		var promises []*Promise
		for i := 0; i < 5; i += 1 {
			p := Issue(ctx, fmt.Sprintf("noop.%d", i))
			promises = append(promises, p)
		}
		if _, err := ReadAll(ctx, promises...); err != nil {
			return err
		}
		return nil
	})

	select {
	case <-allDone():
		if err := ctx.Err(); err != nil {
			t.Error(err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout")
	}
}

func TestReadAll(t *testing.T) {
	ctx, cancel, allDone := NewTestContext()
	defer cancel()

	var promises []*Promise
	for i := 0; i < 5; i += 1 {
		p := Issue(ctx, fmt.Sprintf("noop.%d", i))
		promises = append(promises, p)
	}
	if _, err := ReadAll(ctx, promises...); err != nil {
		t.Error(err)
		cancel()
	}

	select {
	case <-allDone():
		if err := ctx.Err(); err != nil {
			t.Error(err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout")
	}
}

func TestIdempotentRead(t *testing.T) {
	ctx, cancel, allDone := NewTestContext()
	defer cancel()

	Go(ctx, func(ctx context.Context) error {
		p := Issue(ctx, "noop")
		if _, err := Read(ctx, p); err != nil {
			return err
		} else if _, err := Read(ctx, p); err != nil {
			return err
		}
		return nil
	})

	select {
	case <-allDone():
		if err := ctx.Err(); err != nil {
			t.Error(err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout")
	}
}

func TestOneError(t *testing.T) {
	ctx, cancel, allDone := NewTestContext()
	defer cancel()

	myErr := errors.New("myerror")

	Go(ctx, func(ctx context.Context) error {
		Issue(ctx, "noop")
		return myErr
	})

	select {
	case <-allDone():
		if err := ctx.Err(); err != myErr {
			t.Errorf("looking for %q but found %q instead", myErr, err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout")
	}
}

func TestOneErrorOutOfN(t *testing.T) {
	ctx, cancel, allDone := NewTestContext()
	defer cancel()

	myErr := errors.New("myerror")

	for idx := 0; idx < 5; idx += 1 {
		i := idx
		Go(ctx, func(ctx context.Context) error {
			Issue(ctx, "noop")
			if i == 4 {
				return myErr
			} else {
				return nil
			}
		})
	}

	select {
	case <-allDone():
		if err := ctx.Err(); err != myErr {
			t.Errorf("looking for %q but found %q instead", myErr, err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout")
	}
}

func TestCancel(t *testing.T) {
	ctx, cancel, allDone := NewTestContext()
	cancel()

	myErr := context.Canceled
	select {
	case <-allDone():
		if err := ctx.Err(); err != myErr {
			t.Errorf("looking for %q but found %q instead", myErr, err)
		}
	case <-time.After(1 * time.Second):
		t.Fatal("timeout")
	}
}
