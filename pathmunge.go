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

type optType struct {
	del, ps, head bool
} 
var opt optType

func init() {
	flag.BoolVar(&opt.ps, "ps", false, "for Powershell")
	flag.BoolVar(&opt.del, "d", false, "delete the dir from path")
	flag.BoolVar(&opt.head, "h", false, "add to top of list")
}

// Produce PATH environment assignment string
func main() {
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
	plist := genList(orgPath, dir, opt)

	pathListSeparator := string(filepath.ListSeparator)
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

	if opt.ps {
		fmt.Printf("$env:%s=\"%s\"", pathenv, strings.Join(plist, pathListSeparator))
	} else {
		fmt.Printf("%s='%s'\n", pathenv, strings.Join(plist, pathListSeparator))
	}
}

func genList(orgPath, dir string, opt optType) []string {
	plist := []string{}
	in := map[string]struct{}{}
	for _, s := range filepath.SplitList(orgPath) {
		if _, ok := in[s]; !ok && ! (s == dir && (opt.del || opt.head)) {
				plist = append(plist, s)
				in[s] = struct{}{}
		}
	}
	if _, ok := in[dir]; !ok && !opt.del {
		plist = append([]string{dir}, plist...)
	}
	return plist
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
