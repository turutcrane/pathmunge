package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)
// Produce PATH environment assignment string
func main() {
	ps := flag.Bool("ps", false, "for Powershell")
	flag.Parse()
	if flag.NArg() <= 1 {
		log.Fatalf("Usage> %s PATHENV dir\n", filepath.Base(os.Args[0]))
	}
	pathenv := flag.Arg(0)
	dir, err := filepath.Abs(flag.Arg(1))
	if err != nil {
		log.Fatalf("Invalid path specified: %s", flag.Arg(0))
	}

	orgPath := os.Getenv(pathenv)
	plist := []string{}
	in := false
	for _, s := range strings.Split(orgPath, string(os.PathListSeparator)) {
		plist = append(plist, s)
		if s == dir {
			in = true
		}
	}
	if !in {
		plist = append([]string{dir}, plist...)
	}

	pathListSeparator := string(os.PathListSeparator)
	msys, msysRoot := InMsys2()
	if msys {
		for i, p := range plist {
			p = msysPath(p)
			if strings.HasPrefix(p, msysRoot) {
				p = strings.Replace(p, msysRoot, "", 1)
			}
			plist[i] = p
		}
		pathListSeparator = ":"
	}
	for i, p := range plist {
		p = strings.ReplaceAll(p, " ", "\\ ")
		p = strings.ReplaceAll(p, "(", "\\(")
		plist[i] = strings.ReplaceAll(p, ")", "\\)")
	}
	if *ps {
		fmt.Printf("$env:%s=\"%s\"", pathenv, strings.Join(plist, pathListSeparator))
	} else {
		fmt.Printf("%s=%s\n", pathenv, strings.Join(plist, pathListSeparator))
	}
}

// Is it in msys2 environment
func InMsys2() (msys bool, msysRoot string) {
	if os.Getenv("MSYSTEM") != "" {
		msys = true
	}
	if msys {
		// msys2 の / の場所は、
		// $ cygpath -w /
		// の出力で調べられる
		out, err := exec.Command("cygpath", "-w", "/").Output()
		if err != nil {
			log.Fatalf("T62: %v", err)
		}
		msysRoot = strings.Replace(string(out), "\n", "", 1)
		msysRoot = strings.Replace(msysRoot, "\r", "", 1)
		msysRoot = filepath.ToSlash(filepath.Clean(msysRoot))
	}
	return msys, msysPath(msysRoot)
}

// convert C:/ to /C/
func msysPath(p string) string {
	p = filepath.ToSlash(p)
	vol := filepath.VolumeName(p)
	p = strings.Replace(p, vol, "", 1)
	vol = "/" + strings.Replace(vol, ":", "", 1)
	return vol + p
}
