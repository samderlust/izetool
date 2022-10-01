package utils

import (
	"fmt"
	"strings"
)

func StringTemplateReplace(s string, templates map[string]string) string {
	newS := ""

	for key, val := range templates {
		t := fmt.Sprintf("{{%s}}", key)
		newS = strings.ReplaceAll(s, t, val)
	}

	return newS
}

// func processFlagMap() {}
