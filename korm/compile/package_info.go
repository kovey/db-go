package compile

import (
	"encoding/json"
	"os"
	"os/exec"
)

type _packageInfo struct {
	Dir,
	ImportPath,
	Name,
	Target,
	Root,
	StaleReason string
	Stale  bool
	Module struct {
		Main bool
		Path,
		Dir,
		GoMod,
		GoVersion string
	}
	Match,
	GoFiles []string
}

func parsePackageInfo(projectDir, pkgPath string) (*_packageInfo, error) {
	command := []string{"go", "list", "-json", "-find"}
	if pkgPath != "" && pkgPath != "main" {
		command = append(command, pkgPath)
	}
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Dir = projectDir
	cmd.Env = os.Environ()
	bf, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	p := &_packageInfo{}
	err = json.Unmarshal(bf, p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
