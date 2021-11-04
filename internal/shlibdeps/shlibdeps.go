package shlibdeps

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var lddRe = regexp.MustCompile(`^\t([^ ]+) => ([^\s]+)`)

var errLddFailed = errors.New("ldd failed") // sentinel

type LibDep struct {
	Path     string
	Basename string
}

func FindShlibDeps(ldd, fn string, env []string) ([]LibDep, error) {
	cmd := exec.Command(ldd, fn)
	cmd.Env = env
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		// TODO: do not print an error for wrapper programs
		log.Printf("TODO: exclude file %s: %v (out: %s)", fn, err, string(out))
		return nil, errLddFailed // TODO: fix
		return nil, err
	}
	var pkgs []LibDep
	for _, line := range strings.Split(string(out), "\n") {
		matches := lddRe.FindStringSubmatch(line)
		if matches == nil {
			continue
		}
		path, err := filepath.EvalSymlinks(matches[2])
		if err != nil {
			return nil, err
		}
		pkgs = append(pkgs, LibDep{
			Path:     path,
			Basename: filepath.Base(matches[2]),
		})
	}
	return pkgs, nil
}
