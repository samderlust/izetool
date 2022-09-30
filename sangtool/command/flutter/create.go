package flutter

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"com.samderlust/sangtoolbox/sangtool/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	templateFlag = "template"
)

// FlutterCreate create flutter project and customized folders
func FlutterCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create <name>",
		Short: "create a new Flutter project",
		Long:  `create a new Flutter project and template folder`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			cwd, err := os.Getwd()
			if err != nil {
				return errors.Wrap(err, "failed to Getwd")
			}

			template, _ := cmd.Flags().GetString(templateFlag)
			path := filepath.Join(cwd, name)

			templatePath, err := utils.GetTemplateFilePath(template)
			if err != nil {
				return err
			}

			// run flutter create
			createCmd := exec.Command("flutter", "create", name)
			utils.TaskWrapper(
				fmt.Sprintf("Creating Flutter Project: %s ", name),
				func() error {
					return createCmd.Run()
				},
			)

			//read from json template file and create files and folders
			templFile, err := os.ReadFile(templatePath)
			if err != nil {
				return errors.Wrap(err, "failed to read template file")
			}

			var data interface{}
			err = json.Unmarshal(templFile, &data)
			if err != nil {
				return errors.Wrap(err, "reading json err")
			}

			utils.TaskWrapper(
				fmt.Sprintf("Create template: %s ", template),
				func() error {
					return utils.CreateDirsRecursive(data, path)
				},
			)

			return nil
		},
	}

	cmd.Flags().StringP(templateFlag, "t", "example", "the template that will be use, default to example")

	return cmd

}
