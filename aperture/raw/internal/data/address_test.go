package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAddressURIInvalid(t *testing.T) {
	t.Parallel()

	cases := []string{
		// Bad structure
		"",
		"invalid",
		"tcp:",
		// Partial structure
		"tcp:/",
		"tcp://",
		"://host",
		// Complete structure with invalid characters
		"tc3p://host",
		"TCP://host",
		" tcp://localhost:8125",
		"tcp://localhost:8125 ",
		"tcp://loca lhost:8125 ",
	}

	for _, testCase := range cases {
		uri, err := ParseAddressURI(testCase)

		assert.Nil(t, uri)
		assert.Error(t, err)
	}
}

func TestParseAddressURIValid(t *testing.T) {
	t.Parallel()

	cases := []struct {
		raw    string
		parsed *AddressURI
	}{
		{
			raw:    "tcp://localhost:8125",
			parsed: &AddressURI{Scheme: "tcp", Authority: "localhost:8125"},
		},
		{
			raw:    "udp://localhost:8125",
			parsed: &AddressURI{Scheme: "udp", Authority: "localhost:8125"},
		},
		{
			raw:    "unix:///path/to/file",
			parsed: &AddressURI{Scheme: "unix", Authority: "/path/to/file"},
		},
	}

	for _, testCase := range cases {
		uri, err := ParseAddressURI(testCase.raw)

		assert.NoError(t, err)
		assert.Equal(t, testCase.parsed, uri)
	}
}
