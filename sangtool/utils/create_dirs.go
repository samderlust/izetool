package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

// CreateDirsRecursive create folders and files follow provided template
func CreateDirsRecursive(template interface{}, cwd string) error {
	var mapType map[string]interface{}
	var arrType []interface{}

	if reflect.TypeOf(template) != reflect.TypeOf(mapType) {
		return errors.New("invalid template format")
	}

	for key, val := range template.(map[string]interface{}) {
		parentPath := filepath.Join(cwd, key)
		if err := os.MkdirAll(parentPath, 0777); err != nil {
			fmt.Println(err)
		}

		if reflect.TypeOf(val) != reflect.TypeOf(arrType) {
			return errors.New("invalid format")
		}

		for _, v := range val.([]interface{}) {
			childPath := filepath.Join(parentPath, v.(string))
			if strings.Contains(childPath, ".") {
				_, err := os.Create(childPath)
				if err != nil {
					return nil
				}
			} else {
				if err := os.Mkdir(childPath, 0777); err != nil {
					return nil
				}
			}
		}
	}
	return nil
}
