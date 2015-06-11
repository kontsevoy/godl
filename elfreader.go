package main

import (
	"debug/elf"
	"os"
	"path/filepath"
)

func GetELFDependencies(args *Arguments, dlDirs []string) (retval []string) {
	var (
		onFile func(string)
		deps   map[string]bool = make(map[string]bool)
	)

	// gets called on every binary:
	onFile = func(fp string) {
		// no need to check dependencies?
		if args.NoDeps {
			deps[fp] = true
			return
		}
		f, err := elf.Open(fp)
		if err != nil {
			// not an ELF? probaly just a data file:
			fi, err := os.Stat(fp)
			if err == nil && !fi.IsDir() {
				deps[fp] = true
			}
			return
		}
		defer f.Close()

		libs, err := f.ImportedLibraries()
		if err != nil {
			return
		}
		deps[fp] = true

		// check rpath and runpath
		rp1, _ := f.DynString(elf.DT_RPATH)
		rp2, _ := f.DynString(elf.DT_RUNPATH)

		// look for the lib in every location where dynamic linker would look:
		for _, lib := range libs {
			for _, dir := range append(append(rp1, rp2...), dlDirs...) {
				// does this .so file exist? if so, recursively treat it
				// as another binary:
				so := filepath.Join(dir, lib)
				fi, err := os.Stat(so)
				if err == nil && !fi.IsDir() {
					onFile(so)
				}
			}
		}
	}

	// process command-line args (patterns/files):
	for _, p := range GlobMany(args.Patterns, GlobFiles, nil) {
		onFile(p)
	}

	// convert map values to slice of strings:
	retval = make([]string, 0, len(deps)/2)
	for p, _ := range deps {
		retval = append(retval, p)
	}
	return retval
}
