package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gokrazy/freeze/internal/shlibdeps"
)

func copyFile(ctx context.Context, src, dest string) error {
	cp := exec.CommandContext(ctx, "cp", src, dest)
	cp.Stdout = os.Stdout
	cp.Stderr = os.Stderr
	log.Printf("%v", cp.Args)
	if err := cp.Run(); err != nil {
		return fmt.Errorf("%v: %v", cp.Args, err)
	}
	return nil
}

func freeze1(ctx context.Context, wrap, fn string) error {
	log.Printf("%s", fn)
	arg0 := fn
	var args []string
	if wrap != "" {
		arg0 = wrap
		args = []string{fn}
	}
	deps, err := shlibdeps.FindShlibDeps(arg0, args, os.Environ())
	if err != nil {
		return err
	}
	workDir, err := ioutil.TempDir("", "freeze")
	if err != nil {
		return err
	}
	defer func() {
		if err := os.RemoveAll(workDir); err != nil {
			log.Print(err)
		}
	}()

	log.Printf("Copying %s together with its %d ELF shared library dependencies", filepath.Base(fn), len(deps))

	if err := copyFile(ctx, fn, filepath.Join(workDir, filepath.Base(fn))); err != nil {
		return err
	}

	for _, dep := range deps {
		if err := copyFile(ctx, dep.Path, filepath.Join(workDir, dep.Basename)); err != nil {
			return err
		}
	}

	tar := exec.CommandContext(ctx, "tar", "cf", workDir+".tar", filepath.Base(workDir))
	tar.Dir = filepath.Dir(workDir)
	tar.Stdout = os.Stdout
	tar.Stderr = os.Stderr
	log.Printf("%v", tar.Args)
	if err := tar.Run(); err != nil {
		return fmt.Errorf("%v: %v", tar.Args, err)
	}
	log.Printf("Download %s to your gokrazy device and run:\n\tLD_LIBRARY_PATH=$PWD ./ld-linux-x86-64.so.2 ./%s", filepath.Base(workDir+".tar"), filepath.Base(fn))
	// TODO: can we make it a self-extracting archive somehow?
	return nil
}

func freeze() error {
	ctx := context.Background()

	wrap := flag.String("wrap", "", "if non-empty, wrap the command execution (for ldd) with this program, e.g. qemu-aarch64-static when cross-compiling")
	flag.Parse()
	for idx, fn := range flag.Args() {
		if idx > 0 {
			log.Printf("")
		}
		if err := freeze1(ctx, *wrap, fn); err != nil {
			return err
		}
	}

	return nil
}

func main() {
	if err := freeze(); err != nil {
		log.Fatal(err)
	}
}
