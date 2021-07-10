package main

import "github.com/spf13/cobra"

func rootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sangtool",
		Short:   "Tools to help you push up progress",
		Version: "0.0.1",
	}

	cmd.SetVersionTemplate("sangtool CLI v{{.Version}}\n")

	return cmd
}
