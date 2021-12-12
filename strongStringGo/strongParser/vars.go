package strongParser

import (
	"reflect"
	"regexp"
)

var (
	sectionHeader = regexp.MustCompile(`\[([^]]+)\]`)
	keyValue      = regexp.MustCompile(`([^:=\s][^:=]*)\s*(?P<vi>[:=])\s*(.*)$`)
	//continuationLine = regexp.MustCompile(`\w+(.*)$`)
	//interpolater = regexp.MustCompile(`%\(([^)]*)\)s`)
)

var boolMapping = map[string]bool{
	"1":     true,
	"true":  true,
	"on":    true,
	"yes":   true,
	"0":     false,
	"false": false,
	"off":   false,
	"no":    false,
}

var invalidReflectValue = reflect.ValueOf(nil)
