package command

import (
	"com.samderlust/sangtoolbox/sangtool/command/flutter"
	"github.com/spf13/cobra"
)

func FlutterCmd() *cobra.Command {
	flutterCmd := &cobra.Command{
		Use:   "flutter",
		Short: "command lines for Flutter",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				cmd.Help()
			}

			return nil
		},
	}
	flutterCmd.AddCommand(flutter.FlutterCreate())
	flutterCmd.AddCommand(flutter.UploadKeystore())
	return flutterCmd
}
