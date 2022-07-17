package outputer

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

type Output string

const (
	OutputGeneric Output = "generic"
	OutputJSON    Output = "json"

	OutputHelpTextList = `"generic" or "json"`
)

func (o Output) String() string {
	return string(o)
}

func (o Output) Type() string {
	return "output"
}

func (o *Output) Set(v string) error {
	switch v {
	case "":
		*o = OutputGeneric
		return nil
	case "generic", "json":
		*o = Output(v)
		return nil
	default:
		return errors.New(`must be one of "generic" or "json"`)
	}
}

func OutputFlagCompletionFunc(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{
		"generic\tgeneric output",
		"json\toutputs json",
	}, cobra.ShellCompDirectiveDefault
}

func (o *Output) AddFlag(cmd *cobra.Command) error {
	cmd.Flags().VarP(o, "output", "o", "Output format. One of "+OutputHelpTextList)
	return cmd.RegisterFlagCompletionFunc("output", OutputFlagCompletionFunc)
}

type GenericOutputter interface {
	Generic() string
}

func (o Output) Marshal(a any) string {
	switch o {
	case OutputJSON:
		bytes, _ := json.Marshal(a)
		return string(bytes)
	default:
		switch a := a.(type) {
		case string:
			return a
		case GenericOutputter:
			return a.Generic()
		case fmt.Stringer:
			return a.String()
		default:
			return fmt.Sprintf("%v", a)
		}
	}
}
