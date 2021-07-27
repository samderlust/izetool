package utils

import (
	"fmt"
	"log"
	"time"

	"github.com/briandowns/spinner"
)

func TaskWrapper(taskName string, task func() error) error {
	log.SetFlags(0)
	theSpinner := spinner.New(spinner.CharSets[33], 100*time.Millisecond)
	theSpinner.Prefix = fmt.Sprintf("▶️  %s: RUNNING ", taskName)
	theSpinner.FinalMSG = fmt.Sprintf("\n✅  %s: COMPLETED", taskName)
	theSpinner.Start()
	if err := task(); err != nil {
		log.Println(fmt.Sprintf("\n🚫  %s: FAILED", taskName))
		return err
	}

	theSpinner.Stop()
	log.Println("\n---------------------")
	return nil
}
