package command

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

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

			// ctxt := context.Background()
			fmt.Printf("working on %s", name)
			fmt.Printf("template = %s", template)

			//create path
			path := filepath.Join(cwd, name)
			templatePath, _ := filepath.Abs(fmt.Sprintf("./sangtool/templates/%s.json", template))
			fmt.Printf("\npath = %s", path)
			fmt.Printf("\ntemplatePath = %s", templatePath)

			if err := os.Mkdir(path, 0777); err != nil {
				return errors.Wrap(err, fmt.Sprintf("failed to create project %s", name))
			}

			templFile, err := ioutil.ReadFile(templatePath)

			fmt.Printf("file = %s", templFile)

			if err != nil {
				return errors.Wrap(err, "failed to read template file")
			}

			var data interface{}

			err = json.Unmarshal(templFile, &data)

			if err != nil {
				return errors.Wrap(err, "reading json err")
			}

			var arrType []interface{}
			fmt.Printf("\n json = %s", data)
			for k, v := range data.(map[string]interface{}) {
				fmt.Printf("key :: %s === val :: %s", k, v)
				keyPath := (filepath.Join(cwd, fmt.Sprintf("%s/%s", name, k)))
				if err := os.Mkdir(keyPath, 0777); err != nil {
					fmt.Printf("%s exists", name)
				}
				if reflect.TypeOf(v) == reflect.TypeOf(arrType) {
					for _, s := range v.([]interface{}) {

						fmt.Println(s)
						fmt.Println(reflect.TypeOf(s) == reflect.TypeOf(""))
						subPath := filepath.Join(keyPath, s.(string))
						if err := os.Mkdir(subPath, 0777); err != nil {
							fmt.Printf("%s exists", name)
						}

					}
				}
			}

			return nil
		},
	}

	cmd.Flags().StringP(templateFlag, "t", "example", "the tamplate that will be use, default to example")

	return cmd

}
