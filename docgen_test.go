package cobraui

import (
	"reflect"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestConvertflagtypeReturnsRightType(t *testing.T) {

	type test struct {
		input    string
		expected string
	}

	tests := []test{
		{input: "bool", expected: "bool"},
		{input: "boolSlice", expected: "[]bool"},
		{input: "duration", expected: "duration"},
		{input: "int", expected: "int"},
		{input: "intSlice", expected: "[]int"},
		{input: "int32", expected: "int32"},
		{input: "int64", expected: "int64"},
		{input: "string", expected: "string"},
		{input: "stringSlice", expected: "[]string"},
		{input: "stringArray", expected: "stringArray"},
		{input: "blahSlice", expected: "[]blah"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			actual := convertPFlagType(tc.input)
			assert.Equal(t, tc.expected, actual, "expected '%v', got '%v'", tc.expected, actual)
		})
	}

}

func TestCreateOptionDescription(t *testing.T) {
	buildFlag := func(creator func(fs *pflag.FlagSet)) *pflag.Flag {
		fs := &pflag.FlagSet{}
		creator(fs)
		var result *pflag.Flag
		fs.VisitAll(func(f *pflag.Flag) {
			if result == nil {
				result = f
			} else {
				t.Fatal("More than 1 flag found")
			}
		})
		if result == nil {
			t.Fatal("Flag not found")
		}
		return result
	}

	tests := []struct {
		Flag         *pflag.Flag
		Expected     OptionDescription
		ExpectedName string
	}{
		// Booleans
		{
			Flag: buildFlag(func(fs *pflag.FlagSet) { fs.Bool("bool_1", false, "Usage of bool_1") }),
			Expected: OptionDescription{
				Name:        "bool_1",
				Default:     false,
				Description: "Usage of bool_1",
				Type:        "bool",
			},
			ExpectedName: "bool_1",
		},
		{
			Flag: buildFlag(func(fs *pflag.FlagSet) { fs.Bool("bool_2", true, "Usage of bool_2") }),
			Expected: OptionDescription{
				Name:        "bool_2",
				Default:     true,
				Description: "Usage of bool_2",
				Type:        "bool",
			},
			ExpectedName: "bool_2",
		},
		{
			Flag: buildFlag(func(fs *pflag.FlagSet) { fs.BoolP("bool_3", "x", true, "Usage of bool_3") }),
			Expected: OptionDescription{
				Name:        "bool_3",
				Default:     true,
				Description: "Usage of bool_3",
				Type:        "bool",
			},
			ExpectedName: "bool_3",
		},
		// Bool slices
		{
			Flag: buildFlag(func(fs *pflag.FlagSet) { fs.BoolSlice("bs_1", nil, "Usage of bs_1") }),
			Expected: OptionDescription{
				Name:        "bs_1",
				Default:     []bool{},
				Description: "Usage of bs_1",
				Type:        "[]bool",
			},
			ExpectedName: "bs_1",
		},
		{
			Flag: buildFlag(func(fs *pflag.FlagSet) { fs.BoolSlice("bs_2", []bool{true, true, false}, "Usage of bs_2") }),
			Expected: OptionDescription{
				Name:        "bs_2",
				Default:     []bool{true, true, false},
				Description: "Usage of bs_2",
				Type:        "[]bool",
			},
			ExpectedName: "bs_2",
		},
		// Duration
		{
			Flag: buildFlag(func(fs *pflag.FlagSet) { fs.Duration("dur_1", time.Hour, "Usage of dur_1") }),
			Expected: OptionDescription{
				Name:        "dur_1",
				Default:     time.Hour,
				Description: "Usage of dur_1",
				Type:        "duration",
			},
			ExpectedName: "dur_1",
		},
		{
			Flag: buildFlag(func(fs *pflag.FlagSet) { fs.Duration("dur_2", 0, "Usage of dur_2") }),
			Expected: OptionDescription{
				Name:        "dur_2",
				Default:     time.Duration(0),
				Description: "Usage of dur_2",
				Type:        "duration",
			},
			ExpectedName: "dur_2",
		},
		{
			Flag: buildFlag(func(fs *pflag.FlagSet) { fs.DurationP("dur_3", "x", time.Second, "Usage of dur_3") }),
			Expected: OptionDescription{
				Name:        "dur_3",
				Default:     time.Second,
				Description: "Usage of dur_3",
				Type:        "duration",
			},
			ExpectedName: "dur_3",
		},
		// Ints
		{
			Flag: buildFlag(func(fs *pflag.FlagSet) { fs.Int("int_1", 7, "Usage of int_1") }),
			Expected: OptionDescription{
				Name:        "int_1",
				Default:     7,
				Description: "Usage of int_1",
				Type:        "int",
			},
			ExpectedName: "int_1",
		},
		{
			Flag: buildFlag(func(fs *pflag.FlagSet) { fs.Int("int_2", 0, "Usage of int_2") }),
			Expected: OptionDescription{
				Name:        "int_2",
				Default:     0,
				Description: "Usage of int_2",
				Type:        "int",
			},
			ExpectedName: "int_2",
		},
		{
			Flag: buildFlag(func(fs *pflag.FlagSet) { fs.IntP("int_3", "x", 3, "Usage of int_3") }),
			Expected: OptionDescription{
				Name:        "int_3",
				Default:     3,
				Description: "Usage of int_3",
				Type:        "int",
			},
			ExpectedName: "int_3",
		},
		// Int slice
		{
			Flag: buildFlag(func(fs *pflag.FlagSet) { fs.IntSlice("foo", nil, "Usage of foo") }),
			Expected: OptionDescription{
				Name:        "foo",
				Default:     []int{},
				Description: "Usage of foo",
				Type:        "[]int",
			},
			ExpectedName: "foo",
		},
		{
			Flag: buildFlag(func(fs *pflag.FlagSet) { fs.IntSlice("foo", []int{0, 11, 42}, "Usage of foo") }),
			Expected: OptionDescription{
				Name:        "foo",
				Default:     []int{0, 11, 42},
				Description: "Usage of foo",
				Type:        "[]int",
			},
			ExpectedName: "foo",
		},
		// Strings
		{
			Flag: buildFlag(func(fs *pflag.FlagSet) { fs.String("string_1", "", "Usage of string_1") }),
			Expected: OptionDescription{
				Name:        "string_1",
				Default:     "",
				Description: "Usage of string_1",
				Type:        "string",
			},
			ExpectedName: "string_1",
		},
		{
			Flag: buildFlag(func(fs *pflag.FlagSet) { fs.String("string_2", "someValue", "Usage of string_2") }),
			Expected: OptionDescription{
				Name:        "string_2",
				Default:     "someValue",
				Description: "Usage of string_2",
				Type:        "string",
			},
			ExpectedName: "string_2",
		},
		// String slice
		{
			Flag: buildFlag(func(fs *pflag.FlagSet) { fs.StringSlice("foo", nil, "Usage of foo") }),
			Expected: OptionDescription{
				Name:        "foo",
				Default:     []string{},
				Description: "Usage of foo",
				Type:        "[]string",
			},
			ExpectedName: "foo",
		},
		{
			Flag: buildFlag(func(fs *pflag.FlagSet) { fs.StringSlice("foo", []string{"a", "b", "c"}, "Usage of foo") }),
			Expected: OptionDescription{
				Name:        "foo",
				Default:     []string{"a", "b", "c"},
				Description: "Usage of foo",
				Type:        "[]string",
			},
			ExpectedName: "foo",
		},
	}
	for _, test := range tests {
		actual, name, err := createOptionDescription(test.Flag)
		if err != nil {
			t.Errorf("createOptionDescription returned an error: %v", err)
		} else {
			if !reflect.DeepEqual(actual, test.Expected) {
				t.Errorf("createOptionDescription failed. Found '%#v', expected '%#v'", actual, test.Expected)
			}
			if name != test.ExpectedName {
				t.Errorf("createOptionDescription.name failed. Found '%s', expected '%s'", name, test.ExpectedName)
			}
		}
	}
}

func TestGetCommandDetails(t *testing.T) {
	c1 := &cobra.Command{}
	c1.Flags().String("foo", "", "Usage of foo")
	c1.Flags().String("defString", "theDefault", "Usage of defString")
	c1.Flags().String("server.something", "", "Usage of server.something")

	tests := []struct {
		Command  *cobra.Command
		Expected []CommandDetail
	}{
		{
			Command: c1,
			Expected: []CommandDetail{
				{
					Options: OptionDescriptions{
						OptionDescription{
							Name:        "foo",
							Default:     "",
							Description: "Usage of foo",
							Type:        "string",
						},
						OptionDescription{
							Name:        "defString",
							Default:     "theDefault",
							Description: "Usage of defString",
							Type:        "string",
						},
						OptionDescription{
							Name:        "server.something",
							Default:     "theDefault",
							Description: "Usage of server.something",
							Type:        "string",
						},
					},
				},
			},
		},
	}
	for i, test := range tests {
		commandDetails, err := GetCommandDetails(test.Command)
		if err != nil {
			t.Errorf("GetCommandDetails failed on test %d", i)
		} else if reflect.DeepEqual(commandDetails, test.Expected) {
			t.Errorf("Convert failed. Found '%v', expected '%v'", commandDetails, test.Expected)
		}
	}
}
