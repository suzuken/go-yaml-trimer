package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/gobwas/glob"
)

func TestSaveOrder(t *testing.T) {
	in := `T:
  ID:
    type: integer
    format: int32
    x-will-be-removed: true
  Tag:
    type: integer
    format: int32
    x-will-be-removed: false
`
	buf := strings.NewReader(in)

	trimer := &Trimer{}
	err := trimer.open(buf)
	if err != nil {
		t.Fatalf("failed: %s", err)
	}

	var out bytes.Buffer
	err = trimer.Write(&out)
	if err != nil {
		t.Fatalf("failed: %s", err)
	}

	if out.String() != in {
		t.Fatalf("want %s, got %s", in, out.String())
	}

}

func TestTrim(t *testing.T) {
	cases := []struct {
		msg, pattern, in, expected string
	}{
		{
			msg:     "trim x-will-be-removed",
			pattern: `x-will-*`,
			in: `T:
  ID:
    type: integer
    format: int32
    x-will-be-removed: true
  Tag:
    type: integer
    format: int32
    x-will-be-removed: false
`,
			expected: `T:
  ID:
    type: integer
    format: int32
  Tag:
    type: integer
    format: int32
`,
		},
		{
			msg:     "trim multiple",
			pattern: `{x-will-*,format}`,
			in: `T:
  ID:
    type: integer
    format: int32
    x-will-be-removed: true
  Tag:
    type: integer
    format: int32
    x-will-be-removed: false
`,
			expected: `T:
  ID:
    type: integer
  Tag:
    type: integer
`,
		},
	}
	for _, c := range cases {
		t.Run(c.msg, func(t *testing.T) {
			buf := strings.NewReader(c.in)

			trimer := &Trimer{g: glob.MustCompile(c.pattern)}
			err := trimer.open(buf)
			if err != nil {
				t.Fatalf("failed: %s", err)
			}

			var out bytes.Buffer
			err = trimer.Write(&out)
			if err != nil {
				t.Fatalf("failed: %s", err)
			}

			if out.String() != c.expected {
				t.Fatalf("want %s, got %s", c.expected, out.String())
			}
		})
	}
}
