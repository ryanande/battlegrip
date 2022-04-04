package battlegrip

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"

	"github.com/spf13/cobra"
)

// ApplicationDetails is the primary return object.
type ApplicationDetails struct {
	AssemblyName string
	Command      CommandDetail
}

// CommandDetail structure contains parent level commands meta data.
type CommandDetail struct {
	Name             string             `json:"name"`
	Use              string             `json:"use"`
	NameAndAliases   string             `json:"nameandaliases"`
	Aliases          []string           `json:"aliases"`
	Root             string             `json:"root"`
	ShortDescription string             `json:"short"`
	LongDescription  string             `json:"long"`
	Examples         string             `json:"examples"`
	Hidden           bool               `json:"hidden"`
	IsAvailable      bool               `json:"isavailable"`
	HasParent        bool               `json:"hasparent"`
	ParentName       string             `json:"parentname"`
	ParentUse        string             `json:"parentuse"`
	Options          OptionDescriptions `json:"options"`
	Commands         []CommandDetail    `json:"commands"`
}

// OptionDescriptions contains the descriptions for all commandline options of a command.
type OptionDescriptions []OptionDescription

// OptionDescription contains a properties that describe a commandline option.
type OptionDescription struct {
	Name        string      `json:"name"`
	Default     interface{} `json:"default"`
	Description string      `json:"description"`
	Hidden      bool        `json:"hidden"`
	Section     string      `json:"section"`
	Type        string      `json:"type"`
	Values      string      `json:"values,omitempty"`
}

// convertPFlagType converts type names used by pflag to type names understood by the JSON format.
func convertPFlagType(pflagType string) string {
	if strings.HasSuffix(pflagType, "Slice") {
		return "[]" + convertPFlagType(pflagType[:len(pflagType)-len("Slice")])
	}
	switch pflagType {
	default:
		return pflagType
	}
}

type convFunc func(fs *pflag.FlagSet, name string) (interface{}, error)

var (
	convFuncs = map[string]convFunc{
		"bool":        func(fs *pflag.FlagSet, name string) (interface{}, error) { return fs.GetBool(name) },
		"boolSlice":   func(fs *pflag.FlagSet, name string) (interface{}, error) { return fs.GetBoolSlice(name) },
		"duration":    func(fs *pflag.FlagSet, name string) (interface{}, error) { return fs.GetDuration(name) },
		"int":         func(fs *pflag.FlagSet, name string) (interface{}, error) { return fs.GetInt(name) },
		"intSlice":    func(fs *pflag.FlagSet, name string) (interface{}, error) { return fs.GetIntSlice(name) },
		"int32":       func(fs *pflag.FlagSet, name string) (interface{}, error) { return fs.GetInt32(name) },
		"int64":       func(fs *pflag.FlagSet, name string) (interface{}, error) { return fs.GetInt64(name) },
		"uint32":      func(fs *pflag.FlagSet, name string) (interface{}, error) { return fs.GetUint32(name) },
		"uint64":      func(fs *pflag.FlagSet, name string) (interface{}, error) { return fs.GetUint64(name) },
		"string":      func(fs *pflag.FlagSet, name string) (interface{}, error) { return fs.GetString(name) },
		"stringSlice": func(fs *pflag.FlagSet, name string) (interface{}, error) { return fs.GetStringSlice(name) },
		"stringArray": func(fs *pflag.FlagSet, name string) (interface{}, error) { return fs.GetStringArray(name) },
		// Hack: Need to learn if this is the best way to address
		"mapSlice":    func(fs *pflag.FlagSet, name string) (interface{}, error) { return fs.GetStringSlice(name) },
	}
)

// getDefaultValue extracts the default value from the given flag.
func getDefaultValue(flag *pflag.Flag) (interface{}, error) {
	fs := &pflag.FlagSet{}
	fs.AddFlag(flag)
	cf, found := convFuncs[flag.Value.Type()]
	if found {
		if result, err := cf(fs, flag.Name); err == nil {
			return result, nil
		}
		return nil, nil

	}
	return nil, fmt.Errorf("no converter function found for type '%s'", flag.Value.Type())
}

// createOptionDescription creates a desc for the flags.
func createOptionDescription(flag *pflag.Flag) (OptionDescription, string, error) {
	name := flag.Name
	nameParts := strings.Split(name, ".")
	section := ""
	if len(nameParts) > 1 {
		section = nameParts[0]
	}
	defValue, err := getDefaultValue(flag)
	if err != nil {
		return OptionDescription{}, "", err
	}
	d := OptionDescription{
		Name:        name,
		Default:     defValue,
		Description: flag.Usage,
		Section:     section,
		Hidden:      flag.Hidden,
		Type:        convertPFlagType(flag.Value.Type()),
	}
	return d, name, nil
}

// CreateOptionDescriptions maps descriptions for all commandline options.
func createOptionDescriptions(cmd *cobra.Command) (OptionDescriptions, error) {
	result := OptionDescriptions{}

	flags := cmd.Flags()
	var lastErr error
	flags.VisitAll(func(flag *pflag.Flag) {
		d, _, err := createOptionDescription(flag)
		if err != nil {
			lastErr = err
		}
		result = append(result, d)
	})
	if lastErr != nil {
		return nil, lastErr
	}

	return result, nil
}

// GetCommandDetails gathers all details about a command.
func GetCommandDetails(cmd *cobra.Command) (*CommandDetail, error) {
	var destinationCommand CommandDetail
	descriptions, err := createOptionDescriptions(cmd)
	if err != nil {
		fmt.Printf("%v", err)
		return nil, err
	}

	destinationCommand = CommandDetail{
		Name:             cmd.Name(),
		Use:              cmd.Use,
		NameAndAliases:   cmd.NameAndAliases(),
		Aliases:          cmd.Aliases,
		Root:             cmd.Root().Name(),
		ShortDescription: cmd.Short,
		LongDescription:  cmd.Long,
		Examples:         cmd.Example,
		Hidden:           cmd.Hidden,
		IsAvailable:      cmd.IsAvailableCommand(),
		HasParent:        cmd.HasParent(),
		Options:          descriptions,
	}

	if cmd.HasParent() {
		destinationCommand.ParentName = cmd.Parent().Name()
		destinationCommand.ParentUse = cmd.Parent().Use
	}

	for _, c := range cmd.Commands() {
		// skipping not available or help topic command.
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}

		command, err := GetCommandDetails(c)
		if err != nil {
			fmt.Println(err.Error())
		}
		destinationCommand.Commands = append(destinationCommand.Commands, *command)
	}

	return &destinationCommand, nil
}
