package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// GlobMany takes a search pattern and returns absolute file paths that mach that
// pattern.
//	 - targets : list of paths to glob
//   - mask    : GlobDirs or GlobFiles
//   - onErr   : callback function to call when there's an error.
//                can be nil.
func GlobMany(targets []string, onErr func(string, error)) []string {
	rv := make([]string, 0, 20)
	addFile := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			log.Println(err.Error())
			return err
		}
		rv = append(rv, path)
		return err
	}

	for _, p := range targets {
		// "p" is a wildcard pattern? expand it:
		if strings.Contains(p, "*") {
			matches, err := filepath.Glob(p)
			if err == nil {
				// walk each match:
				for _, p := range matches {
					filepath.Walk(p, addFile)
				}
			}
			// path is not a wildcard, walk it:
		} else {
			filepath.Walk(p, addFile)
		}
	}
	return rv
}

// CommonHome takes an array of aboluste file paths and returns a common home
// directory for them
func CommonHome(paths []string) (home string) {
	if len(paths) == 0 {
		return ""
	}
	// first path in list:
	first := paths[0]

	// function returns 'true' if all paths begin with s
	parentForAll := func(s string) bool {
		for _, p := range paths {
			if !strings.HasPrefix(p, s) {
				return false
			}
		}
		return true
	}

	for i := 1; i < len(first); i++ {
		s := first[:i]
		if parentForAll(s) {
			home = s
		} else {
			break
		}
	}
	return home
}
