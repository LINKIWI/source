package data

import (
	"regexp"

	"lib.kevinlin.info/aperture/internal/errors"
)

// AddressURI is a rudimentary representation of a network address presented as a URI.
type AddressURI struct {
	// Scheme is the network associated with the address.
	Scheme string
	// Authority is the address or locator associated with the URI.
	// Note that this field does not strictly conform to the URI format specification.
	Authority string
}

// ParseAddressURI attempts to parse a raw URI into an AddressURI.
func ParseAddressURI(uri string) (*AddressURI, error) {
	addressURIRe := regexp.MustCompile("^(?P<scheme>[a-z]+)://(?P<authority>\\S+)$")

	matches := addressURIRe.FindStringSubmatch(uri)
	if matches == nil || len(matches) != 3 {
		return nil, errors.New(
			"data",
			"address URI format is invalid",
			errors.Tag{Key: "uri", Value: uri},
		)
	}

	return &AddressURI{
		Scheme:    matches[addressURIRe.SubexpIndex("scheme")],
		Authority: matches[addressURIRe.SubexpIndex("authority")],
	}, nil
}
