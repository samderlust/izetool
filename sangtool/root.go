package main

import (
	"com.samderlust/sangtoolbox/sangtool/command"
	"github.com/spf13/cobra"
)

func rootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sangtool",
		Short:   "Tools to help you push up progress",
		Version: "0.0.1",
	}

	cmd.SetVersionTemplate("sangtool CLI v{{.Version}}\n")
	cmd.AddCommand(command.FlutterCmd())

	return cmd
}
