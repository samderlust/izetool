package command

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

//FlutterCreate create flutter project and customized folders
func FlutterCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "flutter create <name>",
		Short: "create a new Flutter project",
		Long:  `create a new Flutter project and template folder`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			cwd, err := os.Getwd()
			if err != nil {
				return errors.Wrap(err, "failed to Getwd")
			}

			// ctxt := context.Background()

			//create path
			path := filepath.Join(cwd, name)
			if err := os.Mkdir(path, 0777); err != nil {
				return errors.Wrap(err, fmt.Sprintf("failed to create project %s", name))
			}

			return nil
		},
	}

	return cmd

}
