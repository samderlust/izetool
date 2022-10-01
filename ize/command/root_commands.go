package command

import (
	"com.samderlust/izetool/ize/command/flutter"
	"com.samderlust/izetool/ize/command/inits"
	"com.samderlust/izetool/ize/command/make"
	"github.com/spf13/cobra"
)

func RootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ize",
		Short:   "Tools to help you push up progress",
		Version: "0.0.1",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				cmd.Help()
			}

			return nil
		},
	}

	cmd.SetVersionTemplate("ize CLI v{{.Version}}\n")
	cmd.AddCommand(inits.InitIzeTool())
	cmd.AddCommand(flutter.FlutterCmd())
	cmd.AddCommand(make.Make())

	return cmd
}
