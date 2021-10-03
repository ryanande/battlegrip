package battlegrip

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"

	"github.com/spf13/cobra"
)

// NewDocs generates the CLI documentation in markdown
func NewJsonDocs(rootCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{
		Use:    "jsondocs",
		Short:  "Generates CLI docs",
		Hidden: true, // this in an internal private command
		Run: func(cmd *cobra.Command, args []string) {

			docs, err := GetCommandDetails(rootCmd)
			if err != nil {
				fmt.Printf("%v", err)
			}

			app := ApplicationDetails{
				AssemblyName: filepath.Base(os.Args[0]),
				Command:      *docs,
			}

			data, err := json.Marshal(app)
			if err != nil {
				fmt.Printf("%v", err)
			}

			dir := "web/src/data"
			if _, err := os.Stat(dir); os.IsNotExist(err) {
				err := os.MkdirAll(dir, 0750)
				if err != nil {
					fmt.Println("Could not create diretory")
				}
			}
			filename := filepath.Join(dir, "commandData.json")
			f, err := os.Create(filename)
			if err != nil {
				fmt.Printf("%v", err)
			}
			defer f.Close()
			if _, err := io.WriteString(f, string(data)); err != nil {
				fmt.Printf("%v", err)
			}
		},
	}

	return cmd
}

type ApplicationDetails struct {
	AssemblyName string
	Command      CommandDetail
}

// CommandDetail structure contains parent level commands meta data
type CommandDetail struct {
	Name             string             `json:"name"`
	IsParent         bool               `json:"isparent"`
	ShortDescription string             `json:"short"`
	LongDescription  string             `json:"long"`
	Examples         string             `json:"examples"`
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
		"string":      func(fs *pflag.FlagSet, name string) (interface{}, error) { return fs.GetString(name) },
		"stringSlice": func(fs *pflag.FlagSet, name string) (interface{}, error) { return fs.GetStringSlice(name) },
		"stringArray": func(fs *pflag.FlagSet, name string) (interface{}, error) { return fs.GetStringArray(name) },
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
		} else {
			return nil, err
		}
	}
	return nil, fmt.Errorf("no converter function found for type '%s'", flag.Value.Type())
}

// createOptionDescription creates a description for the given flag.
// Returns description, name, error.
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

// CreateOptionDescriptions creates a map of descriptions for all the commandline
// options of the given command.
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

// ConvertToJSON converts all the commandline options of the given command to JSON.
func GetCommandDetails(cmd *cobra.Command) (*CommandDetail, error) {
	var destinationCommand CommandDetail
	descriptions, err := createOptionDescriptions(cmd)
	if err != nil {
		fmt.Printf("%v", err)
		return nil, err
	}

	destinationCommand = CommandDetail{
		Name:             cmd.Name(),
		IsParent:         len(cmd.Commands()) > 0,
		ShortDescription: cmd.Short,
		LongDescription:  cmd.Long,
		Examples:         cmd.Example,
		Options:          descriptions,
	}

	for _, c := range cmd.Commands() {
		// skipping not available or help topic command
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
