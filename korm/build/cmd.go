package build

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

type Cmd struct {
	ToolPath   string
	ToolDir    string
	ToolName   string
	ToolArgs   []string
	TempDir    string
	TempGenDir string
	ProjectDir string
}

const (
	version = "0.0.3"
)

func (c *Cmd) Init() {
	c.ToolDir = os.Getenv("GOTOOLDIR")
	c.ToolPath = os.Args[0]
	if len(os.Args) < 2 {
		fmt.Fprintf(flag.CommandLine.Output(), "korm version: %s\n", version)
		os.Exit(0)
	}

	for i, arg := range os.Args[1:] {
		if c.ToolDir != "" && strings.HasPrefix(arg, c.ToolDir) {
			c.ToolName = arg
			if len(os.Args[1:]) > i+1 {
				c.ToolArgs = os.Args[i+2:]
			}
			break
		}
	}

	c.TempDir = path.Join(os.TempDir(), "gobuild_korm_works")
	c.TempGenDir = c.TempDir
	c.ProjectDir, _ = os.Getwd()
	os.MkdirAll(c.TempDir, 0777)
}

var cmd = &Cmd{}
