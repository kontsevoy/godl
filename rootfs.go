package main

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// MakeRootFS creates a working directory where rootfs and manifest
// would reside. If overwrite is set, deletes a pre-existing directory
func MakeRootFS(files []string, args *Arguments) (err error) {
	var (
		srcDir string
		dstDir string
		fi     os.FileInfo
		cutLen int
	)
	// overwrite?
	if args.Force {
		os.RemoveAll(args.RootFS)
	}
	err = os.MkdirAll(args.RootFS, 0771)
	if err != nil {
		return err
	}
	// create manifest and rootfs dir:
	manifest := filepath.Join(args.RootFS, "manifest")
	rootfs := filepath.Join(args.RootFS, "rootfs")
	mbytes := []byte(strings.Replace(DefaultManifest, "%app-name%", args.AppName, 1))
	if args.Manifest != "" {
		mbytes, err = ioutil.ReadFile(args.Manifest)
		if err != nil {
			return err
		}
	}
	ioutil.WriteFile(manifest, mbytes, 0660)
	os.Mkdir(rootfs, 0771)

	// we need to place all found files into a target directory
	// which isn't root. To do that, determine the common parent
	// directory for them on this host, and "chroot" from there
	// into args.Target
	if args.Target != "" {
		cutLen = len(CommonHome(files)) - 1
	}

	// populate rootfs:
	for _, p := range files {
		// create a similar directory under rootfs, with same access flags:
		srcDir = filepath.Dir(p)
		dstDir = filepath.Join(rootfs, args.Target, srcDir[cutLen:])

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
