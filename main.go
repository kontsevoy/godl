package main

import (
	"bufio"
	"debug/elf"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Glob(pattern string, f func(string)) {
	matches, _ := filepath.Glob(pattern)
	for _, fp := range matches {
		fi, err := os.Stat(fp)
		if err == nil && fi.IsDir() == false {
			f(fp)
		}
	}
}

func readDLConfig(pattern string) (dirs []string, err error) {
	f := func(configFile string) {
		fmt.Println("Reading " + configFile)
		fd, err := os.Open(configFile)
		if err != nil {
			return
		}
		defer fd.Close()

		sc := bufio.NewScanner(fd)
		for sc.Scan() {
			line := strings.Trim(sc.Text(), "\t ")
			if len(line) == 0 || line[0] == '#' { // ignore comments and empty lines
				continue
			}

			// found "include" directive?
			words := strings.Fields(line)
			if strings.ToLower(words[0]) == "include" {
				subdirs, err := readDLConfig(words[1])
				if err != nil && !os.IsNotExist(err) {
					return
				}
				dirs = append(dirs, subdirs...)
			}
			dirs = append(dirs, line)
			fmt.Println("\t" + line)
		}
	}

	Glob(pattern, f)
	return dirs, err
}

func getSystemLibdirs() []string {
	dirs, err := readDLConfig("/etc/ld.so.conf")
	if err != nil {
		panic(err)
	}
	return append(dirs, "/usr/lib", "/lib")
}

func main() {
	var f string = '''
	he
	'''
	fmt.Println(f)

	sysDirs := getSystemLibdirs()
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
	Glob("/bin/*", onFile)

	for k, _ := range deps {
		fmt.Println(k)
	}
}
