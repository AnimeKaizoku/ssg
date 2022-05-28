package strongParser

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

// New creates a new ConfigParser.
func NewConfigParser() *ConfigParser {
	return &ConfigParser{
		config:   make(Config),
		defaults: newSection(defaultSectionName),
	}
}

// NewWithDefaults allows creation of a new ConfigParser with a pre-existing
// Dict.
func NewWithDefaults(defaults Dict) *ConfigParser {
	p := ConfigParser{
		config:   make(Config),
		defaults: newSection(defaultSectionName),
	}
	for key, value := range defaults {
		p.defaults.Add(key, value)
	}
	return &p
}

// NewConfigParserFromFile creates a new ConfigParser struct populated from the
// supplied filename.
func NewConfigParserFromFile(filename string) (*ConfigParser, error) {
	p, err := Parse(filename)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// Parse takes a filename and parses it into a ConfigParser value.
func Parse(filename string) (*ConfigParser, error) {
	file, err := os.Open(filename)
	defer func() {
		_ = file.Close()
	}()

	if err != nil {
		return nil, err
	}
	p, err := parseFile(file)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// ParseBytes takes bytes array and parses it into a ConfigParser value.
func ParseBytes(b []byte) (*ConfigParser, error) {
	return parseBytes(b)
}

// ParseBytes takes bytes array and parses it into a ConfigParser value.
func ParseString(value string) (*ConfigParser, error) {
	return parseString(value)
}

func ParseConfig(value any, filename string) error {
	return ParseConfigWithOption(value, filename, nil)
}

func ParseConfigWithOption(value any, filename string, opt *ConfigParserOptions) error {
	p, err := Parse(filename)
	if err != nil {
		if opt == nil && !opt.ReadEnv {
			return err
		}

		p = NewConfigParser()
	}

	p.options = opt
	return parseFinalConfig(value, p)
}

func ParseMainAndArrays[mT any, aT any](filename string, opt *ConfigParserOptions) (*MainAndArrayContainer[mT, aT], error) {
	p, err := Parse(filename)
	if err != nil {
		if opt == nil && !opt.ReadEnv {
			return nil, err
		}

		p = NewConfigParser()
	}

	p.options = opt
	return parseMainAndArrays[mT, aT](p)
}

func ParseMainAndArraysStr[mT any, aT any](valueStr string, opt *ConfigParserOptions) (*MainAndArrayContainer[mT, aT], error) {
	p, err := parseString(valueStr)
	if err != nil {
		if opt == nil && !opt.ReadEnv {
			return nil, err
		}

		p = NewConfigParser()
	}

	p.options = opt
	return parseMainAndArrays[mT, aT](p)
}

func parseMainAndArrays[mT any, aT any](p *ConfigParser) (*MainAndArrayContainer[mT, aT], error) {
	var err error

	if p.options.MainSectionName == "" {
		p.options.MainSectionName = "main"
	}
	container := &MainAndArrayContainer[mT, aT]{
		Main: new(mT),
	}

	for _, current := range p.config {
		if current == nil {
			continue
		} else if current.Name == p.options.MainSectionName {
			err = parseFinalConfigBySection(container.Main, current.Name, p)
			if err != nil {
				return nil, err
			}

			continue
		}

		var currentSection = new(aT)
		err = parseFinalConfigBySection(currentSection, current.Name, p)
		if err != nil {
			return nil, err
		}

		validS, ok := interface{}(currentSection).(SectionValue)
		if ok {
			validS.SetSectionName(current.Name)
		}

		container.Sections = append(container.Sections, currentSection)
	}

	return container, nil
}

func ParseByteConfig(value any, b []byte) error {
	p, err := parseBytes(b)
	if err != nil {
		return err
	}

	return parseFinalConfig(value, p)
}

func ParseStringConfig(value any, strValue string) error {
	p, err := parseString(strValue)
	if err != nil {
		return err
	}

	return parseFinalConfig(value, p)
}

func ParseStringConfigWithOption(value any, strValue string, opt *ConfigParserOptions) error {
	p, err := parseString(strValue)
	if err != nil {
		return err
	}

	p.options = opt

	return parseFinalConfig(value, p)
}

func parseFinalConfigBySection(v any, section string, configValue *ConfigParser) error {
	if configValue.options == nil {
		configValue.options = &ConfigParserOptions{}
	}

	rv := reflect.ValueOf(v)
	myType := reflect.TypeOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidParseError{reflect.TypeOf(v)}
	}

	rv = rv.Elem()
	myType = myType.Elem()

	var currentField reflect.Value
	var shouldSkipCounter bool
	var currentIndex = -1

	for {
		if !shouldSkipCounter {
			currentIndex++
			if currentIndex >= rv.NumField() {
				break
			}
			currentField = rv.Field(currentIndex)
		} else {
			shouldSkipCounter = false
		}

		if !currentField.CanSet() {
			// ignore it
			continue
		}

		switch currentField.Kind() {
		case reflect.Struct:
			continue
		case reflect.Ptr:
			fByName := myType.Field(currentIndex)
			if !fByName.IsExported() {
				continue
			}

			myKind := GetPointerKind(fByName.Type)

			if myKind == reflect.Invalid || myKind == reflect.Struct {
				continue
			}

			SetDefaultValue(currentField, myKind)
			currentField = currentField.Elem()
			shouldSkipCounter = true
			continue
		case reflect.String:
			fByName := myType.Field(currentIndex)
			if !fByName.IsExported() {
				continue
			}

			key := fByName.Tag.Get("key")
			done := false

			theValue, err := configValue.Get(section, key)
			if err == nil {
				// first try: from config file.
				currentField.SetString(theValue)
				continue
			}

			envTag := fByName.Tag.Get("env")
			var envTries []string
			if envTag == "" && configValue.options.ReadEnv {
				// if there is no env tag and we are told to allow
				// reading values from env, try to read it from env.
				if section != "" {
					envTries = append(envTries, strings.ToUpper(section)+"_"+strings.ToUpper(key))
				}
				envTries = append(envTries, strings.ToUpper(key))
			} else {
				// if we are given an env tag, just use that, instead of trying a few times
				// to find the correct variable in env...
				envTries = append(envTries, envTag)
			}

			for _, envTry := range envTries {
				envValue := os.Getenv(envTry)
				if envValue != "" {
					theValue = envValue
					if theValue != "" {
						// second try: read it from env.
						currentField.SetString(theValue)
						break
					}
				}
			}

			if done {
				continue
			}

			theValue = fByName.Tag.Get("default")
			if theValue != "" {
				// third try: from default value.
				currentField.SetString(theValue)
				continue
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fByName := myType.Field(currentIndex)
			if !fByName.IsExported() {
				continue
			}

			key := fByName.Tag.Get("key")
			done := false
			fType := strings.ToLower(fByName.Tag.Get("type"))
			parseValue := func(strValue string) (int64, error) {
				switch fType {
				default:
					// consider int64 for default
					fallthrough
				case "int64":
					return strconv.ParseInt(strValue, 10, 64)
				case "rune":
					return int64(configValue.parseAsRune(strValue)), nil
				}
			}

			theValue, err := configValue.GetIntByType(section, key, fType)
			if err == nil {
				// first try: from config file.
				currentField.SetInt(theValue)
				continue
			}

			envTag := fByName.Tag.Get("env")
			var envTries []string
			if envTag == "" && configValue.options.ReadEnv {
				// if there is no env tag and we are told to allow
				// reading values from env, try to read it from env.
				if section != "" {
					envTries = append(envTries, strings.ToUpper(section)+"_"+strings.ToUpper(key))
				}
				envTries = append(envTries, strings.ToUpper(key))
			} else {
				// if we are given an env tag, just use that, instead of trying a few times
				// to find the correct variable in env...
				envTries = append(envTries, envTag)
			}

			for _, envTry := range envTries {
				envValue := os.Getenv(envTry)
				if envValue != "" {
					theValue, err = parseValue(envValue)
					if err == nil {
						// second try: read it from env.
						currentField.SetInt(theValue)
						break
					}
				}
			}

			if done {
				continue
			}

			theValue, err = parseValue(fByName.Tag.Get("default"))
			if err == nil {
				// third try: from default value.
				currentField.SetInt(theValue)
				continue
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fByName := myType.Field(currentIndex)
			if !fByName.IsExported() {
				continue
			}

			key := fByName.Tag.Get("key")
			done := false

			theValue, err := configValue.GetUint64(section, key)
			if err == nil {
				// first try: from config file.
				currentField.SetUint(uint64(theValue))
				continue
			}

			envTag := fByName.Tag.Get("env")
			var envTries []string
			if envTag == "" && configValue.options.ReadEnv {
				// if there is no env tag and we are told to allow
				// reading values from env, try to read it from env.
				if section != "" {
					envTries = append(envTries, strings.ToUpper(section)+"_"+strings.ToUpper(key))
				}
				envTries = append(envTries, strings.ToUpper(key))
			} else {
				// if we are given an env tag, just use that, instead of trying a few times
				// to find the correct variable in env...
				envTries = append(envTries, envTag)
			}

			for _, envTry := range envTries {
				envValue := os.Getenv(envTry)
				if envValue != "" {
					theValue, err = strconv.ParseUint(envValue, 10, 64)
					if err == nil {
						// second try: read it from env.
						currentField.SetUint(uint64(theValue))
						break
					}
				}
			}

			if done {
				continue
			}

			theValue, err = strconv.ParseUint(fByName.Tag.Get("default"), 10, 64)
			if err == nil {
				// third try: from default value.
				currentField.SetUint(uint64(theValue))
				continue
			}
		case reflect.Bool:
			fByName := myType.Field(currentIndex)
			if !fByName.IsExported() {
				continue
			}

			key := fByName.Tag.Get("key")
			found, done := false, false

			theValue, err := configValue.GetBool(section, key)
			if err == nil {
				// first try: from config file.
				currentField.SetBool(theValue)
				continue
			}

			envTag := fByName.Tag.Get("env")
			var envTries []string
			if envTag == "" && configValue.options.ReadEnv {
				// if there is no env tag and we are told to allow
				// reading values from env, try to read it from env.
				if section != "" {
					envTries = append(envTries, strings.ToUpper(section)+"_"+strings.ToUpper(key))
				}
				envTries = append(envTries, strings.ToUpper(key))
			} else {
				// if we are given an env tag, just use that, instead of trying a few times
				// to find the correct variable in env...
				envTries = append(envTries, envTag)
			}

			for _, envTry := range envTries {
				envValue := os.Getenv(envTry)
				if envValue != "" {
					theValue, found = BoolMapping[envValue]
					if found {
						// second try: read it from env.
						currentField.SetBool(theValue)
						break
					}
				}
			}

			if done {
				continue
			}

			theValue, found = BoolMapping[fByName.Tag.Get("default")]
			if found {
				// third try: from default value.
				currentField.SetBool(theValue)
				continue
			}
		case reflect.Float32, reflect.Float64:
			fByName := myType.Field(currentIndex)
			if !fByName.IsExported() {
				continue
			}

			key := fByName.Tag.Get("key")
			done := false

			theValue, err := configValue.GetFloat64(section, key)
			if err == nil {
				// first try: from config file.
				currentField.SetFloat(theValue)
				continue
			}

			envTag := fByName.Tag.Get("env")
			var envTries []string
			if envTag == "" && configValue.options.ReadEnv {
				// if there is no env tag and we are told to allow
				// reading values from env, try to read it from env.
				if section != "" {
					envTries = append(envTries, strings.ToUpper(section)+"_"+strings.ToUpper(key))
				}
				envTries = append(envTries, strings.ToUpper(key))
			} else {
				// if we are given an env tag, just use that, instead of trying a few times
				// to find the correct variable in env...
				envTries = append(envTries, envTag)
			}

			for _, envTry := range envTries {
				envValue := os.Getenv(envTry)
				if envValue != "" {
					theValue, err = strconv.ParseFloat(envValue, 64)
					if err == nil {
						// second try: read it from env.
						currentField.SetFloat(theValue)
						break
					}
				}
			}

			if done {
				continue
			}

			theValue, err = strconv.ParseFloat(fByName.Tag.Get("default"), 64)
			if err == nil {
				// third try: from default value.
				currentField.SetFloat(theValue)
				continue
			}
		case reflect.Complex64, reflect.Complex128:
			fByName := myType.Field(currentIndex)
			if !fByName.IsExported() {
				continue
			}

			key := fByName.Tag.Get("key")
			done := false

			theValue, err := configValue.GetComplex128(section, key)
			if err == nil {
				// first try: from config file.
				currentField.SetComplex(theValue)
				continue
			}

			envTag := fByName.Tag.Get("env")
			var envTries []string
			if envTag == "" && configValue.options.ReadEnv {
				// if there is no env tag and we are told to allow
				// reading values from env, try to read it from env.
				if section != "" {
					envTries = append(envTries, strings.ToUpper(section)+"_"+strings.ToUpper(key))
				}
				envTries = append(envTries, strings.ToUpper(key))
			} else {
				// if we are given an env tag, just use that, instead of trying a few times
				// to find the correct variable in env...
				envTries = append(envTries, envTag)
			}

			for _, envTry := range envTries {
				envValue := os.Getenv(envTry)
				if envValue != "" {
					theValue, err = strconv.ParseComplex(envValue, 128)
					if err == nil {
						// second try: read it from env.
						currentField.SetComplex(theValue)
						done = true
						break
					}
				}
			}

			if done {
				continue
			}

			theValue, err = strconv.ParseComplex(fByName.Tag.Get("default"), 128)
			if err == nil {
				// third try: from default value.
				currentField.SetComplex(theValue)
				continue
			}
		case reflect.Array, reflect.Slice:
			fByName := myType.Field(currentIndex)
			myKind := getArrayKind(fByName.Type)
			if !fByName.IsExported() {
				continue
			}

			key := fByName.Tag.Get("key")
			fType := strings.ToLower(fByName.Tag.Get("type"))
			isRune := fType == "rune" || fType == "[]rune"

			valueToSet, err := configValue.getArrayValueToSet(section, key, myKind, isRune)
			if err != nil || valueToSet.IsNil() || !valueToSet.IsValid() {
				continue
			}

			currentField.Set(valueToSet)
		}
	}

	return nil
}

func parseFinalConfig(v any, configValue *ConfigParser) error {
	if configValue.options == nil {
		configValue.options = &ConfigParserOptions{}
	}

	rv := reflect.ValueOf(v)
	myType := reflect.TypeOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return &InvalidParseError{reflect.TypeOf(v)}
	}

	rv = rv.Elem()
	myType = myType.Elem()

	var currentField reflect.Value
	var shouldSkipCounter bool
	var currentIndex = -1

	for {
		if !shouldSkipCounter {
			currentIndex++
			if currentIndex >= rv.NumField() {
				break
			}
			currentField = rv.Field(currentIndex)
		} else {
			shouldSkipCounter = false
		}

		if !currentField.CanSet() {
			// ignore it
			continue
		}

		switch currentField.Kind() {
		case reflect.Struct:
			continue
		case reflect.Ptr:
			fByName := myType.Field(currentIndex)
			if !fByName.IsExported() {
				continue
			}

			myKind := GetPointerKind(fByName.Type)

			if myKind == reflect.Invalid || myKind == reflect.Struct {
				continue
			}

			SetDefaultValue(currentField, myKind)
			currentField = currentField.Elem()
			shouldSkipCounter = true
			continue
		case reflect.String:
			fByName := myType.Field(currentIndex)
			if !fByName.IsExported() {
				continue
			}

			section := fByName.Tag.Get("section")
			key := fByName.Tag.Get("key")
			done := false

			theValue, err := configValue.Get(section, key)
			if err == nil {
				// first try: from config file.
				currentField.SetString(theValue)
				continue
			}

			envTag := fByName.Tag.Get("env")
			var envTries []string
			if envTag == "" && configValue.options.ReadEnv {
				// if there is no env tag and we are told to allow
				// reading values from env, try to read it from env.
				if section != "" {
					envTries = append(envTries, strings.ToUpper(section)+"_"+strings.ToUpper(key))
				}
				envTries = append(envTries, strings.ToUpper(key))
			} else {
				// if we are given an env tag, just use that, instead of trying a few times
				// to find the correct variable in env...
				envTries = append(envTries, envTag)
			}

			for _, envTry := range envTries {
				envValue := os.Getenv(envTry)
				if envValue != "" {
					theValue = envValue
					if theValue != "" {
						// second try: read it from env.
						currentField.SetString(theValue)
						break
					}
				}
			}

			if done {
				continue
			}

			theValue = fByName.Tag.Get("default")
			if theValue != "" {
				// third try: from default value.
				currentField.SetString(theValue)
				continue
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fByName := myType.Field(currentIndex)
			if !fByName.IsExported() {
				continue
			}

			section := fByName.Tag.Get("section")
			key := fByName.Tag.Get("key")
			done := false
			fType := strings.ToLower(fByName.Tag.Get("type"))
			parseValue := func(strValue string) (int64, error) {
				switch fType {
				default:
					// consider int64 for default
					fallthrough
				case "int64":
					return strconv.ParseInt(strValue, 10, 64)
				case "rune":
					return int64(configValue.parseAsRune(strValue)), nil
				}
			}

			theValue, err := configValue.GetIntByType(section, key, fType)
			if err == nil {
				// first try: from config file.
				currentField.SetInt(theValue)
				continue
			}

			envTag := fByName.Tag.Get("env")
			var envTries []string
			if envTag == "" && configValue.options.ReadEnv {
				// if there is no env tag and we are told to allow
				// reading values from env, try to read it from env.
				if section != "" {
					envTries = append(envTries, strings.ToUpper(section)+"_"+strings.ToUpper(key))
				}
				envTries = append(envTries, strings.ToUpper(key))
			} else {
				// if we are given an env tag, just use that, instead of trying a few times
				// to find the correct variable in env...
				envTries = append(envTries, envTag)
			}

			for _, envTry := range envTries {
				envValue := os.Getenv(envTry)
				if envValue != "" {
					theValue, err = parseValue(envValue)
					if err == nil {
						// second try: read it from env.
						currentField.SetInt(theValue)
						break
					}
				}
			}

			if done {
				continue
			}

			theValue, err = parseValue(fByName.Tag.Get("default"))
			if err == nil {
				// third try: from default value.
				currentField.SetInt(theValue)
				continue
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			fByName := myType.Field(currentIndex)
			if !fByName.IsExported() {
				continue
			}

			section := fByName.Tag.Get("section")
			key := fByName.Tag.Get("key")
			done := false

			theValue, err := configValue.GetUint64(section, key)
			if err == nil {
				// first try: from config file.
				currentField.SetUint(uint64(theValue))
				continue
			}

			envTag := fByName.Tag.Get("env")
			var envTries []string
			if envTag == "" && configValue.options.ReadEnv {
				// if there is no env tag and we are told to allow
				// reading values from env, try to read it from env.
				if section != "" {
					envTries = append(envTries, strings.ToUpper(section)+"_"+strings.ToUpper(key))
				}
				envTries = append(envTries, strings.ToUpper(key))
			} else {
				// if we are given an env tag, just use that, instead of trying a few times
				// to find the correct variable in env...
				envTries = append(envTries, envTag)
			}

			for _, envTry := range envTries {
				envValue := os.Getenv(envTry)
				if envValue != "" {
					theValue, err = strconv.ParseUint(envValue, 10, 64)
					if err == nil {
						// second try: read it from env.
						currentField.SetUint(uint64(theValue))
						break
					}
				}
			}

			if done {
				continue
			}

			theValue, err = strconv.ParseUint(fByName.Tag.Get("default"), 10, 64)
			if err == nil {
				// third try: from default value.
				currentField.SetUint(uint64(theValue))
				continue
			}
		case reflect.Bool:
			fByName := myType.Field(currentIndex)
			if !fByName.IsExported() {
				continue
			}

			section := fByName.Tag.Get("section")
			key := fByName.Tag.Get("key")
			found, done := false, false

			theValue, err := configValue.GetBool(section, key)
			if err == nil {
				// first try: from config file.
				currentField.SetBool(theValue)
				continue
			}

			envTag := fByName.Tag.Get("env")
			var envTries []string
			if envTag == "" && configValue.options.ReadEnv {
				// if there is no env tag and we are told to allow
				// reading values from env, try to read it from env.
				if section != "" {
					envTries = append(envTries, strings.ToUpper(section)+"_"+strings.ToUpper(key))
				}
				envTries = append(envTries, strings.ToUpper(key))
			} else {
				// if we are given an env tag, just use that, instead of trying a few times
				// to find the correct variable in env...
				envTries = append(envTries, envTag)
			}

			for _, envTry := range envTries {
				envValue := os.Getenv(envTry)
				if envValue != "" {
					theValue, found = BoolMapping[envValue]
					if found {
						// second try: read it from env.
						currentField.SetBool(theValue)
						break
					}
				}
			}

			if done {
				continue
			}

			theValue, found = BoolMapping[fByName.Tag.Get("default")]
			if found {
				// third try: from default value.
				currentField.SetBool(theValue)
				continue
			}
		case reflect.Float32, reflect.Float64:
			fByName := myType.Field(currentIndex)
			if !fByName.IsExported() {
				continue
			}

			section := fByName.Tag.Get("section")
			key := fByName.Tag.Get("key")
			done := false

			theValue, err := configValue.GetFloat64(section, key)
			if err == nil {
				// first try: from config file.
				currentField.SetFloat(theValue)
				continue
			}

			envTag := fByName.Tag.Get("env")
			var envTries []string
			if envTag == "" && configValue.options.ReadEnv {
				// if there is no env tag and we are told to allow
				// reading values from env, try to read it from env.
				if section != "" {
					envTries = append(envTries, strings.ToUpper(section)+"_"+strings.ToUpper(key))
				}
				envTries = append(envTries, strings.ToUpper(key))
			} else {
				// if we are given an env tag, just use that, instead of trying a few times
				// to find the correct variable in env...
				envTries = append(envTries, envTag)
			}

			for _, envTry := range envTries {
				envValue := os.Getenv(envTry)
				if envValue != "" {
					theValue, err = strconv.ParseFloat(envValue, 64)
					if err == nil {
						// second try: read it from env.
						currentField.SetFloat(theValue)
						break
					}
				}
			}

			if done {
				continue
			}

			theValue, err = strconv.ParseFloat(fByName.Tag.Get("default"), 64)
			if err == nil {
				// third try: from default value.
				currentField.SetFloat(theValue)
				continue
			}
		case reflect.Complex64, reflect.Complex128:
			fByName := myType.Field(currentIndex)
			if !fByName.IsExported() {
				continue
			}

			section := fByName.Tag.Get("section")
			key := fByName.Tag.Get("key")
			done := false

			theValue, err := configValue.GetComplex128(section, key)
			if err == nil {
				// first try: from config file.
				currentField.SetComplex(theValue)
				continue
			}

			envTag := fByName.Tag.Get("env")
			var envTries []string
			if envTag == "" && configValue.options.ReadEnv {
				// if there is no env tag and we are told to allow
				// reading values from env, try to read it from env.
				if section != "" {
					envTries = append(envTries, strings.ToUpper(section)+"_"+strings.ToUpper(key))
				}
				envTries = append(envTries, strings.ToUpper(key))
			} else {
				// if we are given an env tag, just use that, instead of trying a few times
				// to find the correct variable in env...
				envTries = append(envTries, envTag)
			}

			for _, envTry := range envTries {
				envValue := os.Getenv(envTry)
				if envValue != "" {
					theValue, err = strconv.ParseComplex(envValue, 128)
					if err == nil {
						// second try: read it from env.
						currentField.SetComplex(theValue)
						done = true
						break
					}
				}
			}

			if done {
				continue
			}

			theValue, err = strconv.ParseComplex(fByName.Tag.Get("default"), 128)
			if err == nil {
				// third try: from default value.
				currentField.SetComplex(theValue)
				continue
			}
		case reflect.Array, reflect.Slice:
			fByName := myType.Field(currentIndex)
			myKind := getArrayKind(fByName.Type)
			if !fByName.IsExported() {
				continue
			}

			section := fByName.Tag.Get("section")
			key := fByName.Tag.Get("key")
			fType := strings.ToLower(fByName.Tag.Get("type"))
			isRune := fType == "rune" || fType == "[]rune"

			valueToSet, err := configValue.getArrayValueToSet(section, key, myKind, isRune)
			if err != nil || valueToSet.IsNil() || !valueToSet.IsValid() {
				continue
			}

			currentField.Set(valueToSet)
		}
	}

	return nil
}

func getArrayKind(t reflect.Type) reflect.Kind {
	myStr := t.String()
	if !strings.HasPrefix(myStr, "[]") {
		return reflect.Invalid
	}

	myStr = myStr[2:]
	myStr = strings.TrimPrefix(myStr, "*")

	return getKind(myStr)
}

func GetPointerKind(t reflect.Type) reflect.Kind {
	myStr := t.String()
	if !strings.HasPrefix(myStr, "*") {
		return reflect.Invalid
	}

	myStr = strings.TrimLeft(myStr, "*")

	return getKind(myStr)
}

func getKind(value string) reflect.Kind {
	switch value {
	case "int":
		return reflect.Int
	case "int8":
		return reflect.Int8
	case "int16":
		return reflect.Int16
	case "int32":
		return reflect.Int32
	case "int64":
		return reflect.Int64
	case "uint":
		return reflect.Uint
	case "uint8":
		return reflect.Uint8
	case "uint16":
		return reflect.Uint16
	case "uint32":
		return reflect.Uint32
	case "uint64":
		return reflect.Uint64
	case "float32":
		return reflect.Float32
	case "float64":
		return reflect.Float64
	case "complex64":
		return reflect.Complex64
	case "complex128":
		return reflect.Complex128
	case "bool":
		return reflect.Bool
	case "string":
		return reflect.String
	}

	return reflect.Invalid
}

func SetDefaultValue(field reflect.Value, kind reflect.Kind) {
	switch kind {
	case reflect.Int:
		var v int
		field.Set(reflect.ValueOf(&v))
	case reflect.Int8:
		var v int8
		field.Set(reflect.ValueOf(&v))
	case reflect.Int16:
		var v int16
		field.Set(reflect.ValueOf(&v))
	case reflect.Int32:
		var v int32
		field.Set(reflect.ValueOf(&v))
	case reflect.Int64:
		var v int64
		field.Set(reflect.ValueOf(&v))
	case reflect.Uint:
		var v uint
		field.Set(reflect.ValueOf(&v))
	case reflect.Uint8:
		var v uint8
		field.Set(reflect.ValueOf(&v))
	case reflect.Uint16:
		var v uint16
		field.Set(reflect.ValueOf(&v))
	case reflect.Uint32:
		var v uint32
		field.Set(reflect.ValueOf(&v))
	case reflect.Uint64:
		var v uint64
		field.Set(reflect.ValueOf(&v))
	case reflect.Float32:
		var v float32
		field.Set(reflect.ValueOf(&v))
	case reflect.Float64:
		var v float64
		field.Set(reflect.ValueOf(&v))
	case reflect.Complex64:
		var v complex64
		field.Set(reflect.ValueOf(&v))
	case reflect.Complex128:
		var v complex128
		field.Set(reflect.ValueOf(&v))
	case reflect.Bool:
		var v bool
		field.Set(reflect.ValueOf(&v))
	case reflect.String:
		var v string
		field.Set(reflect.ValueOf(&v))
	}
}

func GetDefaultValue(kind reflect.Kind) any {
	switch kind {
	case reflect.Int:
		var v int
		return v
	case reflect.Int8:
		var v int8
		return v
	case reflect.Int16:
		var v int16
		return v
	case reflect.Int32:
		var v int32
		return v
	case reflect.Int64:
		var v int64
		return v
	case reflect.Uint:
		var v uint
		return v
	case reflect.Uint8:
		var v uint8
		return v
	case reflect.Uint16:
		var v uint16
		return v
	case reflect.Uint32:
		var v uint32
		return v
	case reflect.Uint64:
		var v uint64
		return v
	case reflect.Float32:
		var v float32
		return v
	case reflect.Float64:
		var v float64
		return v
	case reflect.Complex64:
		var v complex64
		return v
	case reflect.Complex128:
		var v complex128
		return v
	case reflect.Bool:
		var v bool
		return v
	case reflect.String:
		var v string
		return v
	}

	return nil
}

func getNoSectionError(section string) error {
	return fmt.Errorf("no section: '%s'", section)
}

func getNoOptionError(section, option string) error {
	return fmt.Errorf("no option '%s' in section: '%s'", option, section)
}

func parseFile(file *os.File) (*ConfigParser, error) {
	p := NewConfigParser()

	reader := bufio.NewReader(file)
	var lineNo int
	var err error
	var curSect *Section

	for err == nil {
		l, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		lineNo++
		if len(l) == 0 {
			continue
		}
		line := strings.TrimFunc(string(l), unicode.IsSpace)

		// Skip comment lines and empty lines
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		if match := sectionHeader.FindStringSubmatch(line); len(match) > 0 {
			section := match[1]
			if section == defaultSectionName {
				curSect = p.defaults
			} else if _, present := p.config[section]; !present {
				curSect = newSection(section)
				p.config[section] = curSect
			}
		} else if match = keyValue.FindStringSubmatch(line); len(match) > 0 {
			if curSect == nil {
				return nil, fmt.Errorf("missing Section Header: %d %s", lineNo, line)
			}
			curSect.Add(strings.TrimSpace(match[1]), match[3])
		}
	}
	return p, nil
}

func parseBytes(value []byte) (*ConfigParser, error) {
	return parseString(string(value))
}

func parseString(value string) (*ConfigParser, error) {
	p := NewConfigParser()
	allLines := strings.Split(value, "\n")
	var lineNo int
	var curSect *Section

	for _, current := range allLines {

		lineNo++
		if len(current) == 0 {
			continue
		}

		line := strings.TrimSpace(current)

		// Skip comment lines and empty lines
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		if match := sectionHeader.FindStringSubmatch(line); len(match) > 0 {
			section := match[1]
			if section == defaultSectionName {
				curSect = p.defaults
			} else if _, present := p.config[section]; !present {
				curSect = newSection(section)
				p.config[section] = curSect
			}
		} else if match = keyValue.FindStringSubmatch(line); len(match) > 0 {
			if curSect == nil {
				return nil, fmt.Errorf("missing Section Header: %d %s", lineNo, line)
			}

			curSect.Add(strings.TrimSpace(match[1]), match[3])
		}
	}

	return p, nil
}

func parseToStringArray(value string) []string {
	arr := strings.Split(value, ",")

	for i := 0; i < len(arr); i++ {
		arr[i] = strings.TrimSpace(arr[i])
	}

	return arr
}

func parseToIntArray(value string) []int {
	arr := strings.Split(value, ",")
	var myInts []int

	for i := 0; i < len(arr); i++ {
		arr[i] = strings.TrimSpace(arr[i])
		theValue, err := strconv.ParseInt(arr[i], 10, 64)
		if err != nil {
			continue
		}

		myInts = append(myInts, int(theValue))
	}

	return myInts
}

func parseToInt8Array(value string) []int8 {
	arr := strings.Split(value, ",")
	var myInts []int8

	for i := 0; i < len(arr); i++ {
		arr[i] = strings.TrimSpace(arr[i])
		theValue, err := strconv.ParseInt(arr[i], 10, 64)
		if err != nil {
			continue
		}

		myInts = append(myInts, int8(theValue))
	}

	return myInts
}

func parseToInt16Array(value string) []int16 {
	arr := strings.Split(value, ",")
	var myInts []int16

	for i := 0; i < len(arr); i++ {
		arr[i] = strings.TrimSpace(arr[i])
		theValue, err := strconv.ParseInt(arr[i], 10, 64)
		if err != nil {
			continue
		}

		myInts = append(myInts, int16(theValue))
	}

	return myInts
}

func parseToInt32Array(value string, isRune bool) []int32 {
	arr := strings.Split(value, ",")
	var myInts []int32

	for i := 0; i < len(arr); i++ {
		arr[i] = strings.TrimSpace(arr[i])
		if arr[i] == "" {
			continue
		}

		if isRune {
			myInts = append(myInts, int32([]rune(arr[i])[0]))
			continue
		}

		theValue, err := strconv.ParseInt(arr[i], 10, 64)
		if err != nil {
			continue
		}

		myInts = append(myInts, int32(theValue))
	}

	return myInts
}

func parseToInt64Array(value string) []int64 {
	arr := strings.Split(value, ",")
	var myInts []int64

	for i := 0; i < len(arr); i++ {
		arr[i] = strings.TrimSpace(arr[i])
		theValue, err := strconv.ParseInt(arr[i], 10, 64)
		if err != nil {
			continue
		}

		myInts = append(myInts, int64(theValue))
	}

	return myInts
}

func parseToUintArray(value string) []uint {
	arr := strings.Split(value, ",")
	var myInts []uint

	for i := 0; i < len(arr); i++ {
		arr[i] = strings.TrimSpace(arr[i])
		theValue, err := strconv.ParseInt(arr[i], 10, 64)
		if err != nil {
			continue
		}

		myInts = append(myInts, uint(theValue))
	}

	return myInts
}

func parseToUint8Array(value string) []uint8 {
	arr := strings.Split(value, ",")
	var myInts []uint8

	for i := 0; i < len(arr); i++ {
		arr[i] = strings.TrimSpace(arr[i])
		theValue, err := strconv.ParseInt(arr[i], 10, 64)
		if err != nil {
			continue
		}

		myInts = append(myInts, uint8(theValue))
	}

	return myInts
}

func parseToUint16Array(value string) []uint16 {
	arr := strings.Split(value, ",")
	var myInts []uint16

	for i := 0; i < len(arr); i++ {
		arr[i] = strings.TrimSpace(arr[i])
		theValue, err := strconv.ParseInt(arr[i], 10, 64)
		if err != nil {
			continue
		}

		myInts = append(myInts, uint16(theValue))
	}

	return myInts
}

func parseToUint32Array(value string) []uint32 {
	arr := strings.Split(value, ",")
	var myInts []uint32

	for i := 0; i < len(arr); i++ {
		arr[i] = strings.TrimSpace(arr[i])
		theValue, err := strconv.ParseInt(arr[i], 10, 64)
		if err != nil {
			continue
		}

		myInts = append(myInts, uint32(theValue))
	}

	return myInts
}

func parseToUint64Array(value string) []uint64 {
	arr := strings.Split(value, ",")
	var myInts []uint64

	for i := 0; i < len(arr); i++ {
		arr[i] = strings.TrimSpace(arr[i])
		theValue, err := strconv.ParseInt(arr[i], 10, 64)
		if err != nil {
			continue
		}

		myInts = append(myInts, uint64(theValue))
	}

	return myInts
}

func parseToFloat64Array(value string) []float64 {
	arr := strings.Split(value, ",")
	var myFloats []float64

	for i := 0; i < len(arr); i++ {
		arr[i] = strings.TrimSpace(arr[i])
		theValue, err := strconv.ParseFloat(arr[i], 64)
		if err != nil {
			continue
		}

		myFloats = append(myFloats, theValue)
	}

	return myFloats
}

func parseToFloat32Array(value string) []float32 {
	arr := strings.Split(value, ",")
	var myFloats []float32

	for i := 0; i < len(arr); i++ {
		arr[i] = strings.TrimSpace(arr[i])
		theValue, err := strconv.ParseFloat(arr[i], 64)
		if err != nil {
			continue
		}

		myFloats = append(myFloats, float32(theValue))
	}

	return myFloats
}

func parseToComplex64Array(value string) []complex64 {
	arr := strings.Split(value, ",")
	var myFloats []complex64

	for i := 0; i < len(arr); i++ {
		arr[i] = strings.TrimSpace(arr[i])
		theValue, err := strconv.ParseComplex(arr[i], 64)
		if err != nil {
			continue
		}

		myFloats = append(myFloats, complex64(theValue))
	}

	return myFloats
}

func parseToComplex128Array(value string) []complex128 {
	arr := strings.Split(value, ",")
	var myFloats []complex128

	for i := 0; i < len(arr); i++ {
		arr[i] = strings.TrimSpace(arr[i])
		theValue, err := strconv.ParseComplex(arr[i], 128)
		if err != nil {
			continue
		}

		myFloats = append(myFloats, complex128(theValue))
	}

	return myFloats
}

func parseToBoolArray(value string) []bool {
	arr := strings.Split(value, ",")
	var myFloats []bool

	for i := 0; i < len(arr); i++ {
		arr[i] = strings.TrimSpace(arr[i])
		theValue, found := BoolMapping[strings.ToLower(arr[i])]
		if !found {
			continue
		}

		myFloats = append(myFloats, theValue)
	}

	return myFloats
}
