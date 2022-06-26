package migration

import (
	"testing"
)

func TestParse(t *testing.T) {
	t.Run("parse up and down", func(t *testing.T) {
		body := `
### migrate up ###

echo "hi"
rm file
### migrate down ###
# get python version
python3 -v`
		up, down, err := parse(body)
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		if up != `

echo "hi"
rm file
` {
			t.Fatalf("up not expected to equal: %s", up)
		}
		if down != `
# get python version
python3 -v` {
			t.Fatalf("down not expected to equal: %s", down)
		}
	})

	t.Run("wrong order of up and down section", func(t *testing.T) {
		body := `
### migrate down ###
# get python version
python3 -v

### migrate up ###
echo "hi"
rm file
`
		_, _, err := parse(body)
		if err != ErrWrongOrder {
			t.Fatalf("parse expected to fail with ErrWrongOrder, got %v", err)
		}
	})

	t.Run("has only up section", func(t *testing.T) {
		body := `### migrate up ###
echo "hi"`
		up, down, err := parse(body)
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		if up != `
echo "hi"` {
			t.Fatalf("up not expected to equal: %s", up)
		}
		if down != "" {
			t.Fatalf("down expected to be empty, got: %s", down)
		}
	})

	t.Run("has no section definitions", func(t *testing.T) {
		body := `echo "hi"`
		up, down, err := parse(body)
		if err != nil {
			t.Fatalf("err expected to be nil, got: %v", err)
		}
		if up != "" {
			t.Fatalf("up expected to be empty, got: %v", up)
		}
		if down != "" {
			t.Fatalf("down expected to be empty, got: %v", down)
		}
	})

	t.Run("is completely empty", func(t *testing.T) {
		body := ``
		up, down, err := parse(body)
		if err != nil {
			t.Fatalf("err expected to be nil, got: %v", err)
		}
		if up != "" {
			t.Fatalf("up expected to be empty, got: %v", up)
		}
		if down != "" {
			t.Fatalf("down expected to be empty, got: %v", down)
		}
	})
}
