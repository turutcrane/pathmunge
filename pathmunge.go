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
	var msys bool
	if os.Getenv("MSYSTEM") != "" {
		msys = true
	}

	flag.Parse()
	if flag.NArg() == 0 {
		log.Fatalf("Usage> %s path\n", filepath.Base(os.Args[0]))
	}

	path, err := filepath.Abs(flag.Arg(0))
	if err != nil {
		log.Fatalf("Invalid path specified: %s", flag.Arg(0))
	}

	orgPath := os.Getenv("PATH")
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

	pathListSeparator := string(os.PathListSeparator)
	// msys2 の / の場所は、 
	// $ cygpath -w /
	// の出力で調べられる
	if msys {
		for i, p := range(plist) {
			p = filepath.ToSlash(p)
			vol := filepath.VolumeName(p)
			p = strings.Replace(p, vol, "", 1)
			vol = "/" + strings.Replace(vol, ":", "", 1)
			plist[i] = vol + p
		}
		pathListSeparator = ":"
	}
	fmt.Printf("PATH=%s\n", strings.Join(plist, pathListSeparator))
}
