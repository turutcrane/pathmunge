package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	orgPath := os.Getenv("PATH")
	flag.Parse()
	if flag.NArg() == 0 {
		log.Fatalf("Usage> %s path\n", filepath.Base(os.Args[0]))
	}

	path, err := filepath.Abs(flag.Arg(0))
	if err != nil {
		log.Fatalf("Invalid path specified: %s", flag.Arg(0))
	}
	plist := []string{}
	in := false
	for _, s := range strings.Split(orgPath, string(os.PathListSeparator)) {
		plist = append(plist, s)
		if s == path {
			in = true
		}
	}
	if !in {
		plist = append([]string{path}, plist...)
	}
	fmt.Printf("PATH=%s\n", strings.Join(plist, string(os.PathListSeparator)))
}
