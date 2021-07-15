package main

import (
	"flag"
	"fmt"
	"strings"
)

// stringMapFlag implements flag.Value and provides the ability to specify multiple, arbitrary
// key-value parameters, presented as a map of string -> string.
type stringMapFlag struct {
	flag.Value

	values map[string]string
}

// choicesFlag describes a flag.Value that only allows one of several candidate choices.
type choicesFlag struct {
	flag.Value

	candidates    []string
	defaultChoice string
	choice        string
}

func newStringMapFlag() *stringMapFlag {
	return &stringMapFlag{
		values: make(map[string]string),
	}
}

// newChoicesFlag creates a new choicesFlag with candidate choices and a default choice.
func newChoicesFlag(choices []string, defaultChoice string) *choicesFlag {
	return &choicesFlag{candidates: choices, defaultChoice: defaultChoice}
}

func (sm *stringMapFlag) Set(value string) error {
	components := strings.Split(value, "=")
	if len(components) != 2 {
		return fmt.Errorf("improperly formatted key-value parameter")
	}

	sm.values[components[0]] = components[1]

	return nil
}

func (sm *stringMapFlag) Values() map[string]string {
	return sm.values
}

func (sm *stringMapFlag) String() string {
	var kvPairs []string

	for key, value := range sm.values {
		kvPairs = append(kvPairs, fmt.Sprintf("%s=%s", key, value))
	}

	return strings.Join(kvPairs, ",")
}

// Set is called when the associated flag is specified on the CLI.
func (cf *choicesFlag) Set(value string) error {
	for _, candidate := range cf.candidates {
		if value == candidate {
			cf.choice = value
			return nil
		}
	}

	return fmt.Errorf("supplied value is not a valid choice; candidates=%v", cf.candidates)
}

// Choice returns the current choice set by the command-line flag.
func (cf *choicesFlag) Choice() string {
	if cf.choice == "" {
		return cf.defaultChoice
	}

	return cf.choice
}

// String acts as an alias for Choice.
func (cf *choicesFlag) String() string {
	return cf.Choice()
}
