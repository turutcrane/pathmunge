package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Produce PATH environment assignment string
func main() {
	ps := flag.Bool("ps", false, "for Powershell")
	del := flag.Bool("d", false, "delete the dir from path")
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
	for _, s := range filepath.SplitList(orgPath) {
		match := false
		if s == dir {
			match = true
		}
		if !(match && *del) {
			plist = append(plist, s)
		}
		in = match || in
	}
	if !in && !*del {
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
	if runtime.GOOS != "windows" || msys {
		for i, p := range plist {
			p = strings.ReplaceAll(p, " ", "\\ ")
			p = strings.ReplaceAll(p, "(", "\\(")
			plist[i] = strings.ReplaceAll(p, ")", "\\)")
		}
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
