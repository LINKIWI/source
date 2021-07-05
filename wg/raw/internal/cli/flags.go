package cli

import (
	"fmt"
)

// ChoicesFlag describes a flag.Value that only allows one of several candidate choices.
type ChoicesFlag struct {
	candidates    []string
	defaultChoice string
	choice        string
}

// NewChoicesFlag creates a new ChoicesFlag with candidate choices and a default choice.
func NewChoicesFlag(choices []string, defaultChoice string) *ChoicesFlag {
	return &ChoicesFlag{candidates: choices, defaultChoice: defaultChoice}
}

// Set is called when the associated flag is specified on the CLI.
func (cf *ChoicesFlag) Set(value string) error {
	for _, candidate := range cf.candidates {
		if value == candidate {
			cf.choice = value
			return nil
		}
	}

	return fmt.Errorf(
		"flags: supplied value is not a valid choice; candidates=%v",
		cf.candidates,
	)
}

// Choice returns the current choice set by the command-line flag.
func (cf *ChoicesFlag) Choice() string {
	if cf.choice == "" {
		return cf.defaultChoice
	}

	return cf.choice
}

// Candidates returns all available candidates for the flag.
func (cf *ChoicesFlag) Candidates() []string {
	return cf.candidates
}

// String acts as an alias for Choice.
func (cf *ChoicesFlag) String() string {
	return cf.Choice()
}

// ArrayFlag describes a flag.Value with multiple possible values.
type ArrayFlag struct {
	values []string
}

// NewArrayFlag creates a new ArrayFlag.
func NewArrayFlag() *ArrayFlag {
	return &ArrayFlag{values: make([]string, 0)}
}

// Set appends the value to the array in internal state.
func (af *ArrayFlag) Set(value string) error {
	af.values = append(af.values, value)

	return nil
}

// Values returns the values set for the flag.
func (af *ArrayFlag) Values() []string {
	return af.values
}

// String provides a string representation of the string array.
func (af *ArrayFlag) String() string {
	return fmt.Sprintf("%v", af.values)
}
