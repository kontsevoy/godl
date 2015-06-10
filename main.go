package main

import (
	"debug/elf"
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

// Command line arguments:
type Arguments struct {
	Patterns []string
	Force    bool
	DryRun   bool
	Output   string
}

// ParseArgs returns Arguments structure filled with command line arguments.
// If args are invalid, prints help and returns 'false'
func ParseArgs() (bool, *Arguments) {
	cfg := &Arguments{Output: "out.aci"}

	var f = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	f.Usage = func() {
		fmt.Println(UsageStr)
	}
	f.StringVar(&cfg.Output, "o", "out.aci", "")
	f.BoolVar(&cfg.Force, "f", false, "")
	f.BoolVar(&cfg.DryRun, "d", false, "")
	f.Parse(os.Args[1:])

	if len(os.Args) < 2 {
		f.Usage()
		return false, cfg
	}
	cfg.Patterns = f.Args()
	return true, cfg
}

func main() {
	// no arguments?
	canRun, args := ParseArgs()
	if !canRun {
		return
	}

	// for each pattern:
	fl := GlobMany(args.Patterns, GlobFiles, nil)
	for _, p := range fl {
		fmt.Println(p)
	}

	fmt.Println("--------------")

	sysDirs := GetDynLibDirs()
	fmt.Println(sysDirs)

	return

	var (
		onFile func(string)
		deps   map[string]bool = make(map[string]bool)
	)

	onFile = func(fp string) {
		// open
		f, err := elf.Open(fp)
		if err != nil {
			return
		}
		defer f.Close()

		libs, err := f.ImportedLibraries()
		if err != nil {
			fmt.Println("ERROR: " + err.Error())
			return
		}
		deps[fp] = true

		// check rpath and runpath
		rp1, _ := f.DynString(elf.DT_RPATH)
		rp2, _ := f.DynString(elf.DT_RUNPATH)

		// look for the lib in every location where dynamic linker would look:
		for _, lib := range libs {
			for _, dir := range append(append(rp1, rp2...), sysDirs...) {
				libp := filepath.Join(dir, lib)
				fi, err := os.Stat(libp)
				if err == nil && !fi.IsDir() {
					onFile(libp)
				}
			}
		}
	}
}
