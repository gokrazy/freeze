package shlibdeps

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	lddRe = regexp.MustCompile(`^\t(?:[^ ]+) => ([^\s]+)`)
	ldRe  = regexp.MustCompile(`^\t([^ ]+)\s`)
)

var errLddFailed = errors.New("ldd failed") // sentinel

type LibDep struct {
	Path     string
	Basename string
}

func FindShlibDeps(arg0 string, args []string, env []string) ([]LibDep, error) {
	cmd := exec.Command(arg0, args...)
	cmd.Env = append(env, "LD_TRACE_LOADED_OBJECTS=1")
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var pkgs []LibDep
	for _, line := range strings.Split(string(out), "\n") {
		matches := lddRe.FindStringSubmatch(line)
		if matches == nil && strings.Contains(line, "ld-linux") {
			matches = ldRe.FindStringSubmatch(line)
		}
		if matches == nil {
			continue
		}
		if matches[1] == "linux-vdso.so.1" {
			continue // provided by the kernel, not an actual .so file
		}
		path, err := filepath.EvalSymlinks(matches[1])
		if err != nil {
			return nil, err
		}
		pkgs = append(pkgs, LibDep{
			Path:     path,
			Basename: filepath.Base(matches[1]),
		})
	}
	return pkgs, nil
}
