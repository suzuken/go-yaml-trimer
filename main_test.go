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
	err = trimer.write(&out)
	if err != nil {
		t.Fatalf("failed: %s", err)
	}

	if out.String() != in {
		t.Fatalf("want %s, got %s", in, out.String())
	}

}

func TestTrim(t *testing.T) {
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
	expected := `T:
  ID:
    type: integer
    format: int32
  Tag:
    type: integer
    format: int32
`
	buf := strings.NewReader(in)

	trimer := &Trimer{g: glob.MustCompile(`x-will-*`)}
	err := trimer.open(buf)
	if err != nil {
		t.Fatalf("failed: %s", err)
	}

	var out bytes.Buffer
	err = trimer.write(&out)
	if err != nil {
		t.Fatalf("failed: %s", err)
	}

	if out.String() != expected {
		t.Fatalf("want %s, got %s", expected, out.String())
	}

}
