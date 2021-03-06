package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/gobwas/glob"
	yaml "gopkg.in/yaml.v2"
)

type Trimer struct {
	g    glob.Glob
	data yaml.MapSlice
}

func (t *Trimer) open(r io.Reader) error {
	dec := yaml.NewDecoder(r)
	return dec.Decode(&t.data)
}

func (t *Trimer) Open(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return t.open(file)
}

func (t *Trimer) Write(w io.Writer) error {
	t.Trim()

	enc := yaml.NewEncoder(w)
	defer enc.Close()
	return enc.Encode(t.data)
}

func (t *Trimer) write(filename string) error {
	outfile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outfile.Close()
	return t.Write(outfile)
}

func (t *Trimer) trimIter(items yaml.MapSlice) yaml.MapSlice {
	var deleted int
	for i, kv := range items {
		key, ok := kv.Key.(string)
		if !ok {
			continue
		}
		if t.g.Match(key) {
			j := i - deleted
			items = items[:j+copy(items[j:], items[j+1:])]
			deleted++
			continue
		}

		arr, ok := kv.Value.(yaml.MapSlice)
		if !ok {
			continue
		}
		if len(arr) == 0 {
			continue
		}
		items[i].Value = t.trimIter(arr)
	}
	return items
}

func (t *Trimer) Trim() {
	if t.g == nil {
		return
	}

	for i, kv := range t.data {
		key, ok := kv.Key.(string)
		if !ok {
			continue
		}
		if t.g.Match(key) {
			t.data = append(t.data[:i], t.data[i+1:]...)
			continue
		}

		arr, ok := kv.Value.(yaml.MapSlice)
		if !ok {
			continue
		}
		if len(arr) == 0 {
			continue
		}
		t.data[i].Value = t.trimIter(arr)
	}
}

func iferr(err error, msg string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "yaml-trimer: %s, err: %s", msg, err)
		os.Exit(1)
	}
}

func main() {
	var (
		output  = flag.String("output", "", "file name to output. if empty, output to STDOUT")
		pattern = flag.String("pattern", "x-will-*", "glob pattern of YAML property name to remove")
	)
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Fprint(os.Stderr, "Usage: go-yaml-trimer -output output.yaml pattern x-will-* input.yaml\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
	filename := flag.Arg(0)

	g, err := glob.Compile(*pattern)
	iferr(err, "glob compile failed")

	t := &Trimer{g, nil}
	iferr(t.Open(filename), "open YAML file failed")
	if *output == "" {
		iferr(t.Write(os.Stdout), "write YAML file failed")
	} else {
		iferr(t.write(*output), "write YAML file failed")
	}
}
