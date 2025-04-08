package goflags

import (
	"flag"
	"fmt"
	"github.com/cnf/structhash"
	"github.com/google/shlex"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
	"text/tabwriter"
)

// FlagSet is a list of flags for an application
type FlagSet struct {
	CaseSensitive  bool
	Marshal        bool
	description    string
	customHelpText string
	flagKeys       InsertionOrderedMap
	groups         []groupData
	CommandLine    *flag.FlagSet

	// OtherOptionsGroupName is the name for all flags not in a group
	OtherOptionsGroupName string
}

type groupData struct {
	name        string
	description string
}

type FlagData struct {
	usage        string
	short        string
	long         string
	group        string // unused unless set later
	defaultValue interface{}
	skipMarshal  bool
	field        flag.Value
}

// NewFlagSet creates a new flagSet structure for the application
func NewFlagSet() *FlagSet {
	flag.CommandLine.ErrorHandling()
	return &FlagSet{
		flagKeys:              newInsertionOrderedMap(),
		OtherOptionsGroupName: "other options",
		CommandLine:           flag.NewFlagSet(os.Args[0], flag.ExitOnError),
	}
}

func newInsertionOrderedMap() InsertionOrderedMap {
	return InsertionOrderedMap{values: make(map[string]*FlagData)}
}

// SetGroup sets a group with name and description for the command line options
//
// The order in which groups are passed is also kept as is, similar to flags.
func (flagSet *FlagSet) SetGroup(name, description string) {
	flagSet.groups = append(flagSet.groups, groupData{name: name, description: description})
}

// Group sets the group for a flag data
func (flagData *FlagData) Group(name string) {
	flagData.group = name
}

// CreateGroup within the flagset
func (flagSet *FlagSet) CreateGroup(groupName, description string, flags ...*FlagData) {
	flagSet.SetGroup(groupName, description)
	for _, currentFlag := range flags {
		currentFlag.Group(groupName)
	}
}

// StringVarP adds a string flag with a shortname and longname
func (flagSet *FlagSet) StringVarP(field *string, long, short, defaultValue, usage string) *FlagData {
	flagData := &FlagData{
		usage:        usage,
		long:         long,
		defaultValue: defaultValue,
	}
	if short != "" {
		flagData.short = short
		flagSet.CommandLine.StringVar(field, short, defaultValue, usage)
		flagSet.flagKeys.Set(short, flagData)
	}
	flagSet.CommandLine.StringVar(field, long, defaultValue, usage)
	flagSet.flagKeys.Set(long, flagData)
	return flagData
}

// StringVar adds a string flag with a longname
func (flagSet *FlagSet) StringVar(field *string, long, defaultValue, usage string) *FlagData {
	return flagSet.StringVarP(field, long, "", defaultValue, usage)
}

// BoolVarP adds a bool flag with a shortname and longname
func (flagSet *FlagSet) BoolVarP(field *bool, long, short string, defaultValue bool, usage string) *FlagData {
	flagData := &FlagData{
		usage:        usage,
		long:         long,
		defaultValue: strconv.FormatBool(defaultValue),
	}
	if short != "" {
		flagData.short = short
		flagSet.CommandLine.BoolVar(field, short, defaultValue, usage)
		flagSet.flagKeys.Set(short, flagData)
	}
	flagSet.CommandLine.BoolVar(field, long, defaultValue, usage)
	flagSet.flagKeys.Set(long, flagData)
	return flagData
}

// BoolVar adds a bool flag with a longname
func (flagSet *FlagSet) BoolVar(field *bool, long string, defaultValue bool, usage string) *FlagData {
	return flagSet.BoolVarP(field, long, "", defaultValue, usage)
}

// IntVarP adds a int flag with a shortname and longname
func (flagSet *FlagSet) IntVarP(field *int, long, short string, defaultValue int, usage string) *FlagData {
	flagData := &FlagData{
		usage:        usage,
		short:        short,
		long:         long,
		defaultValue: strconv.Itoa(defaultValue),
	}
	if short != "" {
		flagData.short = short
		flagSet.CommandLine.IntVar(field, short, defaultValue, usage)
		flagSet.flagKeys.Set(short, flagData)
	}
	flagSet.CommandLine.IntVar(field, long, defaultValue, usage)
	flagSet.flagKeys.Set(long, flagData)
	return flagData
}

// IntVar adds a int flag with a longname
func (flagSet *FlagSet) IntVar(field *int, long string, defaultValue int, usage string) *FlagData {
	return flagSet.IntVarP(field, long, "", defaultValue, usage)
}

// Parse parses the flags provided to the library.
func (flagSet *FlagSet) Parse(args ...string) error {
	flagSet.CommandLine.SetOutput(os.Stdout)
	flagSet.CommandLine.Usage = flagSet.usageFunc
	toParse := os.Args[1:]
	if len(args) > 0 {
		toParse = args
	}
	_ = flagSet.CommandLine.Parse(toParse)

	return nil
}

func (flagSet *FlagSet) usageFunc() {
	var helpAsked bool

	// Only show help usage if asked by user
	for _, arg := range os.Args {
		argStripped := strings.Trim(arg, "-")
		if argStripped == "h" || argStripped == "help" {
			helpAsked = true
		}
	}
	if !helpAsked {
		return
	}

	cliOutput := flagSet.CommandLine.Output()
	fmt.Fprintf(cliOutput, "%s\n\n", flagSet.description)
	fmt.Fprintf(cliOutput, "Usage:\n  %s [flags]\n\n", os.Args[0])
	fmt.Fprintf(cliOutput, "Flags:\n")

	writer := tabwriter.NewWriter(cliOutput, 0, 0, 1, ' ', 0)

	// If a user has specified a group with help, and we have groups, return with the tool's usage function
	if len(flagSet.groups) > 0 && len(os.Args) == 3 {
		group := flagSet.getGroupbyName(strings.ToLower(os.Args[2]))
		if group.name != "" {
			flagSet.displayGroupUsageFunc(newUniqueDeduper(), group, cliOutput, writer)
			return
		}
		flag := flagSet.getFlagByName(os.Args[2])
		if flag != nil {
			flagSet.displaySingleFlagUsageFunc(os.Args[2], flag, cliOutput, writer)
			return
		}
	}

	if len(flagSet.groups) > 0 {
		flagSet.usageFuncForGroups(cliOutput, writer)
	} else {
		flagSet.usageFuncInternal(writer)
	}

	// If there is a custom help text specified, print it
	if !(strings.TrimSpace(flagSet.customHelpText) == "") {
		fmt.Fprintf(cliOutput, "\n%s\n", flagSet.customHelpText)
	}

}

func (flagSet *FlagSet) getGroupbyName(name string) groupData {
	for _, group := range flagSet.groups {
		if strings.EqualFold(group.name, name) || strings.EqualFold(group.description, name) {
			return group
		}
	}
	return groupData{}
}

func (flagSet *FlagSet) getFlagByName(name string) *FlagData {
	var flagData *FlagData
	flagSet.flagKeys.forEach(func(key string, data *FlagData) {
		// check if the items are equal
		// - Case sensitive
		equal := flagSet.CaseSensitive && (data.long == name || data.short == name)
		// - Case insensitive
		equalFold := !flagSet.CaseSensitive && (strings.EqualFold(data.long, name) || strings.EqualFold(data.short, name))
		if equal || equalFold {
			flagData = data
			return
		}
	})
	return flagData
}

// usageFuncForGroups prints usage for command line flags with grouping enabled
func (flagSet *FlagSet) usageFuncForGroups(cliOutput io.Writer, writer *tabwriter.Writer) {
	uniqueDeduper := newUniqueDeduper()

	var otherOptions []string
	for _, group := range flagSet.groups {
		otherOptions = append(otherOptions, flagSet.displayGroupUsageFunc(uniqueDeduper, group, cliOutput, writer)...)
	}

	// Print Any additional flag that may have been left
	if len(otherOptions) > 0 {
		fmt.Fprintf(cliOutput, "%s:\n", normalizeGroupDescription(flagSet.OtherOptionsGroupName))

		for _, option := range otherOptions {
			fmt.Fprint(writer, option, "\n")
		}
		writer.Flush()
	}
}

// displayGroupUsageFunc displays usage for a group
func (flagSet *FlagSet) displayGroupUsageFunc(uniqueDeduper *uniqueDeduper, group groupData, cliOutput io.Writer, writer *tabwriter.Writer) []string {
	fmt.Fprintf(cliOutput, "%s:\n", normalizeGroupDescription(group.description))

	var otherOptions []string
	flagSet.flagKeys.forEach(func(key string, data *FlagData) {
		if currentFlag := flagSet.CommandLine.Lookup(key); currentFlag != nil {
			if data.group == "" {
				if !uniqueDeduper.isUnique(data) {
					return
				}
				otherOptions = append(otherOptions, createUsageString(data, currentFlag))
				return
			}
			// Ignore the flag if it's not in our intended group
			if !strings.EqualFold(data.group, group.name) {
				return
			}
			if !uniqueDeduper.isUnique(data) {
				return
			}
			result := createUsageString(data, currentFlag)
			fmt.Fprint(writer, result, "\n")
		}
	})
	writer.Flush()
	fmt.Printf("\n")
	return otherOptions
}

// displaySingleFlagUsageFunc displays usage for a single flag
func (flagSet *FlagSet) displaySingleFlagUsageFunc(name string, data *FlagData, cliOutput io.Writer, writer *tabwriter.Writer) {
	if currentFlag := flagSet.CommandLine.Lookup(name); currentFlag != nil {
		result := createUsageString(data, currentFlag)
		fmt.Fprint(writer, result, "\n")
		writer.Flush()
	}
}

type uniqueDeduper struct {
	hashes map[string]interface{}
}

func newUniqueDeduper() *uniqueDeduper {
	return &uniqueDeduper{hashes: make(map[string]interface{})}
}

// usageFuncInternal prints usage for command line flags
func (flagSet *FlagSet) usageFuncInternal(writer *tabwriter.Writer) {
	uniqueDeduper := newUniqueDeduper()

	flagSet.flagKeys.forEach(func(key string, data *FlagData) {
		if currentFlag := flagSet.CommandLine.Lookup(key); currentFlag != nil {
			if !uniqueDeduper.isUnique(data) {
				return
			}
			result := createUsageString(data, currentFlag)
			fmt.Fprint(writer, result, "\n")
		}
	})
	writer.Flush()
}

// isUnique returns true if the flag is unique during iteration
func (u *uniqueDeduper) isUnique(data *FlagData) bool {
	dataHash := data.Hash()
	if _, ok := u.hashes[dataHash]; ok {
		return false // Don't print the value if printed previously
	}
	u.hashes[dataHash] = struct{}{}
	return true
}

// Hash returns the unique hash for a flagData structure
// NOTE: Hash panics when the structure cannot be hashed.
func (flagData *FlagData) Hash() string {
	hash, _ := structhash.Hash(flagData, 1)
	return hash
}

func createUsageString(data *FlagData, currentFlag *flag.Flag) string {
	valueType := reflect.TypeOf(currentFlag.Value)

	result := createUsageFlagNames(data)
	result += createUsageTypeAndDescription(currentFlag, valueType)
	result += createUsageDefaultValue(data, currentFlag, valueType)

	return result
}

func createUsageDefaultValue(data *FlagData, currentFlag *flag.Flag, valueType reflect.Type) string {
	if !isZeroValue(currentFlag, currentFlag.DefValue) {
		defaultValueTemplate := " (default "
		switch valueType.String() { // ugly hack because "flag.stringValue" is not exported from the parent library
		case "*flag.stringValue":
			defaultValueTemplate += "%q"
		default:
			defaultValueTemplate += "%v"
		}
		defaultValueTemplate += ")"
		return fmt.Sprintf(defaultValueTemplate, data.defaultValue)
	}
	return ""
}

// isZeroValue determines whether the string represents the zero
// value for a flag.
func isZeroValue(f *flag.Flag, value string) bool {
	// Build a zero value of the flag's Value type, and see if the
	// result of calling its String method equals the value passed in.
	// This works unless the Value type is itself an interface type.
	valueType := reflect.TypeOf(f.Value)
	var zeroValue reflect.Value
	if valueType.Kind() == reflect.Ptr {
		zeroValue = reflect.New(valueType.Elem())
	} else {
		zeroValue = reflect.Zero(valueType)
	}
	return value == zeroValue.Interface().(flag.Value).String()
}

func createUsageTypeAndDescription(currentFlag *flag.Flag, valueType reflect.Type) string {
	var result string

	flagDisplayType, usage := flag.UnquoteUsage(currentFlag)
	if len(flagDisplayType) > 0 {
		if flagDisplayType == "value" { // hardcoded in the goflags library
			switch valueType.Kind() {
			case reflect.Ptr:
				pointerTypeElement := valueType.Elem()
				switch pointerTypeElement.Kind() {
				case reflect.Slice, reflect.Array:
					switch pointerTypeElement.Elem().Kind() {
					case reflect.String:
						flagDisplayType = "string[]"
					default:
						flagDisplayType = "value[]"
					}
				}
			}
		}
		result += " " + flagDisplayType
	}

	result += "\t\t"
	result += strings.ReplaceAll(usage, "\n", "\n"+strings.Repeat(" ", 4)+"\t")
	return result
}

func isEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

func createUsageFlagNames(data *FlagData) string {
	flagNames := strings.Repeat(" ", 2) + "\t"

	var validFlags []string
	addValidParam := func(value string) {
		if !isEmpty(value) {
			validFlags = append(validFlags, fmt.Sprintf("-%s", value))
		}
	}

	addValidParam(data.short)
	addValidParam(data.long)

	if len(validFlags) == 0 {
		panic("CLI arguments cannot be empty.")
	}

	flagNames += strings.Join(validFlags, ", ")
	return flagNames
}

// normalizeGroupDescription returns normalized description field for group
func normalizeGroupDescription(description string) string {
	return strings.ToUpper(description)
}

// GetArgsFromString allows splitting a string into arguments
// following the same rules as the shell.
func GetArgsFromString(str string) []string {
	args, _ := shlex.Split(str)
	return args
}
