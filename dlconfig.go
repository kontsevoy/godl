package main

import (
	"bufio"
	"os"
	"strings"
)

// GetDynLibDirs returns locations (directories) where the dynamic linker (DL)
// looks for .so libraries when it launches applications.
//
// This information is read from /etc/ld.so.conf file
func GetDynLibDirs() []string {
	dirs, err := ParseDynLibConf("/etc/ld.so.conf")
	if err != nil {
		panic(err)
	}
	return append(dirs, "/usr/lib", "/lib")
}

// ParseDynLibConf reads/parses DL config files defined as a pattern
// and returns a list of directories found in there (or an error).
func ParseDynLibConf(pattern string) (dirs []string, err error) {
	files := GlobMany([]string{pattern}, nil)
	for _, configFile := range files {
		fd, err := os.Open(configFile)
		if err != nil {
			return dirs, err
		}
		defer fd.Close()

		sc := bufio.NewScanner(fd)
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			// ignore comments and empty lines
			if len(line) == 0 || line[0] == '#' || line[0] == ';' {
				continue
			}
			// found "include" directive?
			words := strings.Fields(line)
			if strings.ToLower(words[0]) == "include" {
				subdirs, err := ParseDynLibConf(words[1])
				if err != nil && !os.IsNotExist(err) {
					return dirs, err
				}
				dirs = append(dirs, subdirs...)
			} else {
				dirs = append(dirs, line)
			}
		}
	}
	return dirs, err
}
