package main

import (
	"os"
	"path/filepath"
)

const (
	GlobDirs = 1 << iota
	GlobFiles
)

// GlobMany takes a search pattern and returns files that mach that
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
