package build

import (
	"log"
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

func (c *Cmd) Init() {
	c.ToolDir = os.Getenv("GOTOOLDIR")
	c.ToolPath = os.Args[0]
	if c.ToolDir == "" {
		log.Println("GOTOOLDIR is empty")
	}

	if len(os.Args) < 2 {
		log.Println("args error")
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
	if err := os.MkdirAll(c.TempDir, 0777); err != nil {
		log.Println("Init() fail, os.MkdirAll tempDir", err)
	}
}

var cmd = &Cmd{}
