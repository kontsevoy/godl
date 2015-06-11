package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func MakeACI(dir string, out string) error {
	out, _ = filepath.Abs(out)

	// can we create an ACI file with such name?
	f, err := os.Create(out)
	if err != nil {
		return err
	}
	f.Close()

	o, err := exec.Command("tar", "czf", out, "-C", dir, "manifest", "rootfs").Output()
	if err != nil {
		log.Println(o)
		return err
	}
	return nil
}
