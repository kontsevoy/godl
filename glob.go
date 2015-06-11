package main

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	GlobDirs = 1 << iota
	GlobFiles
)

const (
	MaxPath = 255 // equivalent to MAX_PATH in c
)

// GlobMany takes a search pattern and returns absolute file paths that mach that
// pattern.
//	 - pattenrs : list of strings like "/usr/**/*" similar to filepath.Glob
//   - mask     : GlobDirs or GlobFiles
//   - onErr    : callback function to call when there's an error.
//                can be nil.
func GlobMany(patterns []string, mask int, onErr func(string, error)) []string {
	rv := make([]string, 0, 20)
	addFile := func(f string) {
		rv = append(rv, f)
	}
	for _, p := range patterns {
		matches, _ := filepath.Glob(p)
		for _, fp := range matches {
			fp, _ = filepath.Abs(fp)
			fi, err := os.Stat(fp)
			if err != nil && onErr != nil {
				onErr(fp, err)
				continue
			}
			// dir?
			if (mask&GlobDirs == GlobDirs) && fi.IsDir() {
				addFile(fp)
				// file?
			} else if (mask&GlobFiles == GlobFiles) && !fi.IsDir() {
				addFile(fp)
			}
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
