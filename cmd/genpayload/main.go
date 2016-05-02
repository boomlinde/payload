package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/boomlinde/payload"
)

func getWalker(p payload.Payload, root string) filepath.WalkFunc {
	return func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !f.IsDir() {
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}

			fmt.Fprintf(os.Stderr, "Loading file \"%s\"\n", rel)
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			p[rel] = data
		}
		return nil
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <path> [path ...]\n", os.Args[0])
		os.Exit(-1)
	}

	p := make(payload.Payload)

	for _, path := range os.Args[1:] {
		s, err := os.Lstat(path)
		if err != nil {
			panic(err)
		}
		if s.IsDir() {
			if err := filepath.Walk(os.Args[1], getWalker(p, path)); err != nil {
				panic(err)
			}
		} else {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				panic(err)
			}
			p[s.Name()] = data
		}
	}
	p.Dump(os.Stdout)
	fmt.Fprintf(os.Stderr, "Done\n")
}
