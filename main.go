package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Command line arguments:
type Arguments struct {
	Patterns []string
	Force    bool
	NoDeps   bool
	RootFS   string
	OutACI   string
	Target   string
	Manifest string
	AppName  string
	AppDesc  string
}

// ParseArgs returns Arguments structure filled with command line arguments.
// If args are invalid, prints help and returns 'false'
func ParseArgs() (bool, *Arguments) {
	cfg := &Arguments{}

	var f = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	f.Usage = func() {
		fmt.Println(UsageStr)
	}
	f.StringVar(&cfg.OutACI, "o", "", "output rootfs location")
	f.StringVar(&cfg.RootFS, "r", "aci", "output aci image")
	f.BoolVar(&cfg.Force, "f", false, "overwrite existing output")
	f.BoolVar(&cfg.NoDeps, "i", false, "ignore dependencies")
	f.StringVar(&cfg.Target, "t", "", "target directory within rootfs")
	f.StringVar(&cfg.Manifest, "m", "", "manifest file")
	f.StringVar(&cfg.AppName, "n", "", "application name")
	f.StringVar(&cfg.AppDesc, "d", "", "application description")
	f.Parse(os.Args[1:])

	if len(f.Args()) == 0 {
		f.Usage()
		return false, cfg
	}
	cfg.Patterns = f.Args()

	// no manifest specified? make sure we generate an app name (if not provided)
	if cfg.Manifest == "" && cfg.AppName == "" {
		if cfg.OutACI != "" {
			cfg.AppName = strings.Split(cfg.OutACI, ".")[0]
		} else {
			cfg.AppName = filepath.Base(cfg.Patterns[0])
		}
	}
	return true, cfg
}

func main() {
	log.SetFlags(0)

	// no arguments?
	canContinue, args := ParseArgs()
	if !canContinue {
		return
	}

	// get a list of dependencies (binaries and .so libs) from
	// the user-provided patterns:
	deps := GetELFDependencies(args, GetDynLibDirs())
	for _, f := range deps {
		fmt.Println("Adding " + f)
	}

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
