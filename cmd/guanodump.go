package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/riggsd/guano-go/guano"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s FILE...\n", filepath.Base(os.Args[0]))
		os.Exit(2)
	}
	for _, fname := range os.Args[1:] {
		fmt.Printf("%s\n", fname)
		g, err := guano.Read(fname)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed reading %q: %v\n\n", fname, err)
			continue
		}
		for k, v := range g.Fields {
			fmt.Printf("%s:\t%s\n", k, v)
		}
		fmt.Println()
	}
}
