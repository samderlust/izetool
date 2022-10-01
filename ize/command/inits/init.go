package inits

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"com.samderlust/izetool/ize/utils"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

func InitIzeTool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "create templates folder",
		Long:  "create templates folder to contain the boiled-plate templates",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			utils.LogDone("initializing")
			theSpinner := spinner.New(spinner.CharSets[33], 100*time.Millisecond)
			theSpinner.Start()

			home, err := os.UserHomeDir()
			if err != nil {
				log.Fatal(err)
			}

			toolDir := filepath.Join(home, "sangtool_templates")
			os.Mkdir(toolDir, 0777)

			exPath := filepath.Join(toolDir, "example.json")

			_, b, _, _ := runtime.Caller(0)
			basePath := filepath.Join(filepath.Dir(b), "../..")
			templatePath := filepath.Join(basePath, "templates/example.json")

			if err := utils.CopyFile(templatePath, exPath); err != nil {
				log.Fatal(err)
			}
			theSpinner.Stop()

			utils.LogDone(fmt.Sprintf("template created at %s", toolDir))

			return nil
		},
	}

	return cmd
}
