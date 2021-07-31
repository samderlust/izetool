package utils

import (
	"fmt"
	"log"
)

// LogDone log when a task is done
func LogDone(msg string) {
	log.Println(fmt.Sprintf("\n✅  %s", msg))
}

// LogStart log when a task is done
func LogStart(msg string) {
	log.Println(fmt.Sprintf("\n▶️  %s", msg))
}
