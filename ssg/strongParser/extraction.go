package strongParser

import (
	"os"
	"reflect"
	"strconv"
	"strings"
)

// extractFieldValue
func extractFieldValue[T comparable](
	parser *ConfigParser,
	myType reflect.Type,
	currentIndex int, section string,
	converter fieldValueConverter[T]) T {

	var realDefault T
	var resultValue T
	fByName := myType.Field(currentIndex)
	if section == "" {
		section = fByName.Tag.Get("section")
	}

	if section == "" && parser.options.MainSectionName != "" {
		section = parser.options.MainSectionName
	}

	key := fByName.Tag.Get("key")
	if key == "" {
		// convert the field name to snake case
		key = toSnakeCase(fByName.Name)
	}

	fType := strings.ToLower(fByName.Tag.Get("type"))
	theValue, err := parser.Get(section, key)
	if err == nil {
		// first try: from config file.
		resultValue, err = converter(fType, theValue)
		if err == nil && resultValue != realDefault {
			return resultValue
		}
	}

	envTag := fByName.Tag.Get("env")
	var envTries []string
	if envTag == "" && parser.options.ReadEnv {
		// if there is no env tag and we are told to allow
		// reading values from env, try to read it from env.
		if section != "" {
			envTries = append(envTries, strings.ToUpper(section)+"_"+strings.ToUpper(key))
		}
		envTries = append(envTries, key)
		envTries = append(envTries, strings.ToUpper(key))
	} else {
		// if we are given an env tag, just use that, instead of trying a few times
		// to find the correct variable in env...
		envTries = append(envTries, envTag)
	}

	for _, envTry := range envTries {
		envValue := os.Getenv(envTry)
		if envValue != "" {
			resultValue, err = converter(fType, envValue)
			if err == nil && resultValue != realDefault {
				return resultValue
			}
		}
	}

	resultValue, _ = converter(fType, fByName.Tag.Get("default"))
	return resultValue
}

func toSnakeCase(s string) string {
	var result []rune
	for i, c := range s {
		if i > 0 && c >= 'A' && c <= 'Z' {
			result = append(result, '_')
		}
		// lower-case the c and append
		result = append(result, c|0x20)
	}
	return string(result)
}

func extractStr(fType, s string) (string, error) { return s, nil }

func parseAsRune(value string) rune {
	if value == "" {
		return 0
	}

	return ([]rune(value))[0]
}

func extractInt64(fType, strValue string) (int64, error) {
	switch fType {
	default:
		// consider int64 for default
		fallthrough
	case "int64":
		return strconv.ParseInt(strValue, 10, 64)
	case "rune":
		return int64(parseAsRune(strValue)), nil
	}
}

func extractUInt64(_, strValue string) (uint64, error) {
	return strconv.ParseUint(strValue, 10, 64)
}

func extractBool(fType, strValue string) (bool, error) {
	return BoolMapping[strings.ToLower(strValue)], nil
}

func extractFloat64(fType, strValue string) (float64, error) {
	return strconv.ParseFloat(strValue, 64)
}

func extractComplex128(fType, strValue string) (complex128, error) {
	return strconv.ParseComplex(strValue, 128)
}
