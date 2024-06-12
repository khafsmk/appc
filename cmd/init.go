package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/khafsmk/appc/tmpl"
	"github.com/spf13/cobra"
)

var (
	initCmd = &cobra.Command{
		Use:     "init",
		Aliases: []string{"create"},
		Short:   "Initialize a web project",
		Run: func(cmd *cobra.Command, args []string) {
			_, err := initProject(args)
			cobra.CheckErr(err)
		},
	}
)

func initProject(args []string) (string, error) {
	if len(args) < 1 {
		return "", fmt.Errorf("directory is required")
	}
	base := args[0]
	if err := os.MkdirAll(base, os.ModePerm); err != nil {
		return "", fmt.Errorf("create directory %q failed: %w", base, err)
	}
	if err := tmpl.Gen(base); err != nil {
		return "", fmt.Errorf("generate project failed: %w", err)
	}
	return path.Join(base, base), nil
}
