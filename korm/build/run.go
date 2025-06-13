package build

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/kovey/db-go/korm/compile"
)

func Run() {
	cmd.Init()
	if cmd.ToolName == "" {
		os.Exit(2)
	}
	toolName := filepath.Base(cmd.ToolName)
	switch strings.TrimSuffix(toolName, ".exe") {
	case "compile":
		compile.Compile(cmd.ProjectDir, cmd.TempDir, cmd.ToolArgs)
	}

	cmd := exec.Command(cmd.ToolName, cmd.ToolArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	if err := cmd.Run(); err != nil {
		log.Println(err)
	}
}
