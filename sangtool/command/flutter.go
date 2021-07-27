package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"com.samderlust/sangtoolbox/sangtool/utils"
	"github.com/briandowns/spinner"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	templateFlag = "template"
)

//FlutterCreate create flutter project and customized folders
func FlutterCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "flutter_create <name>",
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

			_, b, _, _ := runtime.Caller(0)
			basepath := filepath.Join(filepath.Dir(b), "../..")
			templatePath := filepath.Join(basepath, fmt.Sprintf("sangtool/templates/%s.json", template))

			createCmd := exec.Command("flutter", "create", name)

			createSpinner := spinner.New(spinner.CharSets[35], 100*time.Millisecond)
			createSpinner.Prefix = (fmt.Sprintf("Creating Flutter Project: %s ", name))
			createSpinner.Start()

			if err := createCmd.Run(); err != nil {
				fmt.Println(err)
				fmt.Println("failed creating Flutter project")
				createSpinner.Stop()
				return err
			}
			createSpinner.Stop()

			fmt.Println("\nFlutter project created!!")

			templFile, err := ioutil.ReadFile(templatePath)

			if err != nil {
				return errors.Wrap(err, "failed to read template file")
			}

			var data interface{}
			err = json.Unmarshal(templFile, &data)

			if err != nil {
				return errors.Wrap(err, "reading json err")
			}

			tempSpinner := spinner.New(spinner.CharSets[35], 100*time.Millisecond)

			tempSpinner.Prefix = (fmt.Sprintf("Creating template: %s ", template))
			tempSpinner.Start()

			if err := utils.CreateDirsRecursive(data, path); err != nil {
				fmt.Println("\nFailed creating template")
				return err
			}

			tempSpinner.Stop()
			fmt.Println("\nCreating template finished")
			return nil
		},
	}

	cmd.Flags().StringP(templateFlag, "t", "example", "the tamplate that will be use, default to example")

	return cmd

}
