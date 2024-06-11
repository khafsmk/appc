package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	initCmd = &cobra.Command{
		Use:     "init",
		Aliases: []string{"create"},
		Short:   "Initialize a web project",
		Run: func(cmd *cobra.Command, args []string) {
			path, err := initProject(args)
			cobra.CheckErr(err)
			fmt.Printf("Project initialized at %s\n", path)
		},
	}
)

func initProject(args []string) (string, error) {
	return "", nil
}
