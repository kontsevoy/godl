package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// Command line arguments:
type Arguments struct {
	Patterns []string
	Force    bool
	RootFS   string
	OutACI   string
	Target   string
	Manifest string
}

// ParseArgs returns Arguments structure filled with command line arguments.
// If args are invalid, prints help and returns 'false'
func ParseArgs() (bool, *Arguments) {
	cfg := &Arguments{}

	var f = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	f.Usage = func() {
		fmt.Println(UsageStr)
	}
	f.StringVar(&cfg.OutACI, "o", "out.aci", "")
	f.StringVar(&cfg.RootFS, "r", "aci", "")
	f.BoolVar(&cfg.Force, "f", false, "")
	f.StringVar(&cfg.Target, "t", "", "")
	f.StringVar(&cfg.Manifest, "m", "", "")
	f.Parse(os.Args[1:])

	if len(os.Args) < 2 {
		f.Usage()
		return false, cfg
	}
	cfg.Patterns = f.Args()
	return true, cfg
}

func main() {
	log.SetFlags(0)

	// no arguments?
	canRun, args := ParseArgs()
	if !canRun {
		return
	}

	// get a list of dependencies (binaries and .so libs) from
	// the user-provided patterns:
	deps := GetELFDependencies(args.Patterns, GetDynLibDirs())

	// create a directory which will hold rootfs+manifest
	err := MakeRootFS(deps, args)
	if err != nil {
		log.Fatal(err)
	}

	// make an ACI (if there was an option for it)
	if args.OutACI != "" {
		err = MakeACI(args.RootFS, args.OutACI)
		if err != nil {
			log.Fatal(err)
		}
	}

	// success:
	fmt.Printf("Created RootFS in %v\n", args.RootFS)
	if args.OutACI != "" {
		fmt.Printf("Created ACI image: %v", args.OutACI)
	}
	fmt.Println("")
}
