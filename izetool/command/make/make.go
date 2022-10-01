package make

import (
	"encoding/json"
	"fmt"
	"os"

	"com.samderlust/izetool/izetool/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	nameFlag = "name"
)

func Make() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "make <template>",
		Short: "make files and folders",
		Long:  "make files and folder recursively with provided template",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			template := args[0]

			cwd, err := os.Getwd()
			if err != nil {
				return errors.Wrap(err, "failed to Getwd")
			}

			templatePath, err := utils.GetTemplateFilePath(template)
			if err != nil {
				return errors.Wrap(err, "failed to template file")
			}

			templFile, err := os.ReadFile(templatePath)
			if err != nil {
				return errors.Wrap(err, "failed to read template file")
			}

			var data interface{}
			err = json.Unmarshal(templFile, &data)
			if err != nil {
				return errors.Wrap(err, "reading json err")
			}

			theName, _ := cmd.Flags().GetString(nameFlag)

			utils.TaskWrapper(
				fmt.Sprintf("Create template: %s ", template),
				func() error {
					return utils.CreateDirsRecursiveWithName(data, map[string]string{nameFlag: theName}, cwd)
				},
			)

			return nil
		},
	}
	cmd.Flags().StringP(nameFlag, "n", "example", "the template that will be use, default to example")

	return cmd
}
