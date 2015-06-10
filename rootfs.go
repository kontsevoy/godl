package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func MakeRootFS(files []string, binarypath string) (err error) {
	var (
		tempDir string
		srcDir  string
		dstDir  string
		fi      os.FileInfo
	)

	// prepare directory structure to keep our ACI content
	if false {
		tempDir, err := ioutil.TempDir("", "aci")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tempDir)
		// debug:
	} else {
		tempDir = "/tmp/aci"
		os.RemoveAll(tempDir)
		os.Mkdir(tempDir, 0771)
	}

	manifest := filepath.Join(tempDir, "manifest")
	rootfs := filepath.Join(tempDir, "rootfs")

	// create manifest and rootfs dir:
	ioutil.WriteFile(manifest, []byte(DefaultManifest), 0660)
	os.Mkdir(rootfs, 0771)

	// create rootfs:
	for _, p := range files {
		// create a similar directory under rootfs, with same access flags:
		srcDir = filepath.Dir(p)
		dstDir = filepath.Join(rootfs, srcDir)
		fi, err = os.Stat(srcDir)
		if err != nil {
			return err
		}
		os.MkdirAll(dstDir, fi.Mode())

		// try to create a hard link first (much faster):
		dest := filepath.Join(dstDir, filepath.Base(p))
		err = os.Link(p, dest)
		if err != nil {
			// fall back to copying a file:
			if err = CopyFile(p, dest); err != nil {
				return err
			}
		}
	}
	return nil
}

// CopyFile makes a copy of a file preserving its access flags
func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	fi, err := in.Stat()
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, fi.Mode())
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	err = out.Sync()
	return nil
}
