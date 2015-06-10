package main

import (
	"flag"
	"fmt"
	"os"
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

	// get a list of dependencies (binaries and .so libs) from
	// the user-provided patterns:
	deps := GetELFDependencies(args.Patterns, GetDynLibDirs())

	// look at what we've got!
	for _, p := range deps {
		fmt.Println(p)
	}
}
