package strongParser

import (
	"errors"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

func (p *ConfigParser) isDefaultSection(section string) bool {
	return section == defaultSectionName
}

// Defaults returns the items in the map used for default values.
func (p *ConfigParser) Defaults() Dict {
	return p.defaults.Items()
}

// Sections returns a list of section names, excluding [DEFAULT].
func (p *ConfigParser) Sections() []string {
	sections := make([]string, 0)
	for section := range p.config {
		sections = append(sections, section)
	}
	sort.Strings(sections)
	return sections
}

// HasSection returns true if the named section is present in the
// configuration.
//
// The DEFAULT section is not acknowledged.
func (p *ConfigParser) HasSection(section string) bool {
	_, present := p.config[section]
	return present
}

// Options returns a list of option names for the given section name.
//
// Returns an error if the section does not exist.
func (p *ConfigParser) Options(section string) ([]string, error) {
	if !p.HasSection(section) {
		return nil, getNoSectionError(section)
	}

	seenOptions := make(map[string]bool)
	for _, option := range p.config[section].Options() {
		seenOptions[option] = true
	}

	for _, option := range p.defaults.Options() {
		seenOptions[option] = true
	}

	options := make([]string, 0)
	for option := range seenOptions {
		options = append(options, option)
	}

	sort.Strings(options)

	return options, nil
}

// Get returns string value for the named option.
//
// Returns an error if a section does not exist
// Returns an error if the option does not exist either in the section or in
// the defaults
func (p *ConfigParser) Get(section, option string) (string, error) {
	if section == "" || option == "" {
		return "", errors.New("section and option must be non-empty")
	}

	if !p.HasSection(section) {
		if !p.isDefaultSection(section) {
			return "", getNoSectionError(section)
		}
		if value, err := p.defaults.Get(option); err != nil {
			return "", getNoOptionError(section, option)
		} else {
			return value, nil
		}
	} else if value, err := p.config[section].Get(option); err == nil {
		return value, nil
	} else if value, err := p.defaults.Get(option); err == nil {
		return value, nil
	}

	return "", getNoOptionError(section, option)
}

func (p *ConfigParser) GetIntByType(section, option, fType string) (int64, error) {
	result, err := p.Get(section, option)
	if err != nil {
		return 0, err
	}

	var value int64
	switch fType {
	default:
		// consider int64 for default
		fallthrough
	case "int64", "int32", "int", "int16", "int8":
		value, err = strconv.ParseInt(result, 10, 64)
	case "rune":
		value = int64(p.parseAsRune(result))
	}

	return value, err
}

// ItemsWithDefaults returns a copy of the named section Dict including
// any values from the Defaults.
//
// NOTE: This is different from the Python version which returns a list of
// tuples
func (p *ConfigParser) ItemsWithDefaults(section string) (Dict, error) {
	if !p.HasSection(section) {
		return nil, getNoSectionError(section)
	}
	s := make(Dict)

	for k, v := range p.defaults.Items() {
		s[k] = v
	}

	for k, v := range p.config[section].Items() {
		s[k] = v
	}

	return s, nil
}

// Items returns a copy of the section Dict not including the Defaults.
//
// NOTE: This is different from the Python version which returns a list of
// tuples
func (p *ConfigParser) Items(section string) (Dict, error) {
	if !p.HasSection(section) {
		return nil, getNoSectionError(section)
	}

	return p.config[section].Items(), nil
}

// Set puts the given option into the named section.
//
// Returns an error if the section does not exist.
func (p *ConfigParser) Set(section, option, value string) error {
	var setSection *Section

	if p.isDefaultSection(section) {
		setSection = p.defaults
	} else if _, present := p.config[section]; !present {
		return getNoSectionError(section)
	} else {
		setSection = p.config[section]
	}
	setSection.Add(option, value)
	return nil
}

func (p *ConfigParser) GetInt64(section, option string) (int64, error) {
	result, err := p.Get(section, option)
	if err != nil {
		return 0, err
	}

	value, err := strconv.ParseInt(result, 10, 64)
	if err != nil {
		return 0, err
	}

	return value, nil
}

func (p *ConfigParser) GetUint64(section, option string) (uint64, error) {
	result, err := p.Get(section, option)
	if err != nil {
		return 0, err
	}

	value, err := strconv.ParseUint(result, 10, 64)
	if err != nil {
		return 0, err
	}

	return value, nil
}

func (p *ConfigParser) GetRune(section, option string) rune {
	result, err := p.Get(section, option)
	if err != nil {
		return 0
	}

	if result == "" {
		return 0
	}

	return ([]rune(result))[0]
}

func (p *ConfigParser) parseAsRune(value string) rune {
	return ([]rune(value))[0]
}

func (p *ConfigParser) GetFloat64(section, option string) (float64, error) {
	result, err := p.Get(section, option)
	if err != nil {
		return 0, err
	}

	value, err := strconv.ParseFloat(result, 64)
	if err != nil {
		return 0, err
	}

	return value, nil
}

func (p *ConfigParser) GetComplex128(section, option string) (complex128, error) {
	result, err := p.Get(section, option)
	if err != nil {
		return 0, err
	}

	value, err := strconv.ParseComplex(result, 128)
	if err != nil {
		return 0i, err
	}

	return value, nil
}

func (p *ConfigParser) GetBool(section, option string) (bool, error) {
	result, err := p.Get(section, option)
	if err != nil {
		return false, err
	}

	result = strings.ToLower(result)

	booleanValue, present := BoolMapping[result]
	if !present {
		return false, fmt.Errorf("not a boolean: '%s'", result)
	}

	return booleanValue, nil
}

func (p *ConfigParser) GetStringSlice(section, option string) ([]string, error) {
	result, err := p.Get(section, option)
	if err != nil {
		return nil, err
	}

	return parseToStringArray(result), nil
}

// getArrayValueToSet returns array value to set.
// s is section; o is option; k is kind.
func (p *ConfigParser) getArrayValueToSet(s, o string, k reflect.Kind, isRune bool) (rValue, error) {
	result, err := p.Get(s, o)
	if err != nil {
		return invalidReflectValue, err
	}

	switch k {
	case reflect.String:
		return reflect.ValueOf(parseToStringArray(result)), nil
	case reflect.Int:
		return reflect.ValueOf(parseToIntArray(result)), nil
	case reflect.Int8:
		return reflect.ValueOf(parseToInt8Array(result)), nil
	case reflect.Int16:
		return reflect.ValueOf(parseToInt16Array(result)), nil
	case reflect.Int32:
		return reflect.ValueOf(parseToInt32Array(result, isRune)), nil
	case reflect.Int64:
		return reflect.ValueOf(parseToInt64Array(result)), nil
	case reflect.Uint:
		return reflect.ValueOf(parseToUintArray(result)), nil
	case reflect.Uint8:
		return reflect.ValueOf(parseToUint8Array(result)), nil
	case reflect.Uint16:
		return reflect.ValueOf(parseToUint16Array(result)), nil
	case reflect.Uint32:
		return reflect.ValueOf(parseToUint32Array(result)), nil
	case reflect.Uint64:
		return reflect.ValueOf(parseToUint64Array(result)), nil
	case reflect.Float32:
		return reflect.ValueOf(parseToFloat32Array(result)), nil
	case reflect.Float64:
		return reflect.ValueOf(parseToFloat64Array(result)), nil
	case reflect.Complex64:
		return reflect.ValueOf(parseToComplex64Array(result)), nil
	case reflect.Complex128:
		return reflect.ValueOf(parseToComplex128Array(result)), nil
	case reflect.Bool:
		return reflect.ValueOf(parseToBoolArray(result)), nil
	}

	return invalidReflectValue, fmt.Errorf("unsupported kind: %s", k.String())
}

func (p *ConfigParser) GetIntSlice(section, option string) ([]int64, error) {
	result, err := p.Get(section, option)
	if err != nil {
		return nil, err
	}

	return parseToInt64Array(result), nil
}

func (p *ConfigParser) HasOption(section, option string) (bool, error) {
	var s *Section
	if p.isDefaultSection(section) {
		s = p.defaults
	} else if _, present := p.config[section]; !present {
		return false, getNoSectionError(section)
	} else {
		s = p.config[section]
	}

	_, err := s.Get(option)
	return err == nil, nil
}

func (p *ConfigParser) RemoveOption(section, option string) error {
	var s *Section
	if p.isDefaultSection(section) {
		s = p.defaults
	} else if _, present := p.config[section]; !present {
		return getNoSectionError(section)
	} else {
		s = p.config[section]
	}
	return s.Remove(option)
}

//---------------------------------------------------------

func (s *Section) Add(key, value string) error {
	lookupKey := s.safeKey(key)

	s.options[key] = s.safeValue(value)
	s.lookup[lookupKey] = key
	return nil
}

func (s *Section) Get(key string) (string, error) {
	lookupKey, present := s.lookup[s.safeKey(key)]
	if !present {
		return "", getNoOptionError(s.Name, key)
	}

	if value, present := s.options[lookupKey]; present {
		return strings.TrimSpace(value), nil
	}

	return "", getNoOptionError(s.Name, key)
}

func (s *Section) Options() []string {
	return s.options.Keys()
}

func (s *Section) Items() Dict {
	return s.options
}

func (s *Section) safeValue(in string) string {
	return strings.TrimSpace(in)
}

func (s *Section) safeKey(in string) string {
	return strings.ToLower(strings.TrimSpace(in))
}

func (s *Section) Remove(key string) error {
	_, present := s.options[key]
	if !present {
		return getNoOptionError(s.Name, key)
	}

	// delete doesn't return anything, but this does require
	// that the passed key to be removed matches the options key.
	delete(s.lookup, s.safeKey(key))
	delete(s.options, key)
	return nil
}

func newSection(name string) *Section {
	return &Section{
		Name:    name,
		options: make(Dict),
		lookup:  make(Dict),
	}
}

//---------------------------------------------------------

// Keys returns a sorted slice of keys
func (d Dict) Keys() []string {
	var keys []string

	for key := range d {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	return keys
}

//---------------------------------------------------------

func NewChainMap(dicts ...Dict) *ChainMap {
	chainMap := &ChainMap{
		maps: make([]Dict, 0),
	}

	chainMap.maps = append(chainMap.maps, dicts...)

	return chainMap
}

func (c *ChainMap) Len() int {
	return len(c.maps)
}

func (c *ChainMap) Get(key string) string {
	var value string

	for _, dict := range c.maps {
		if result, present := dict[key]; present {
			value = result
		}
	}
	return value
}

//---------------------------------------------------------

func (e *InvalidParseError) Error() string {
	if e.Type == nil {
		return "strongParser: Parse(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "strongParser: Parse(non-pointer " + e.Type.String() + ")"
	}

	return "strongParser: Parse(nil " + e.Type.String() + ")"
}

//---------------------------------------------------------
