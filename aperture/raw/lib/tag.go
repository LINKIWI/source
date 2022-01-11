package lib

import (
	"fmt"
	"strings"
)

// item is a key-value pair.
type item struct {
	key   string
	value interface{}
}

// Field is a list of zero or more key-value pairs from which a full or partial set of serialized
// tags may be derived.
type Field struct {
	items []item
}

// String returns a human-readable representation of the current list of items in the Field.
func (f Field) String() string {
	components := make([]string, len(f.items))

	for i, it := range f.items {
		components[i] = fmt.Sprintf("[%s: %v]", it.key, it.value)
	}

	return strings.Join(components, " ")
}

// Tag creates a new Field directly from a tag key and value.
func Tag(key string, value interface{}) Field {
	return Field{[]item{{key, value}}}
}

// MapTag creates a new Field with zero or more items from the contents of an existing map.
func MapTag(tags map[string]interface{}) Field {
	if tags == nil {
		return Field{[]item{}}
	}

	items := make([]item, len(tags))

	i := 0
	for key, value := range tags {
		items[i] = item{key, value}
		i++
	}

	return Field{items}
}

// CombineTags serializes a list of Fields into a map that is suitable for use as the `tags`
// parameter for aperture.Statsd emission methods. In the case of tag key conflicts, precedence is
// determined by order; tags that appear later in the list with an existing key name are given
// precedence of value over those that appear earlier in the list.
//
// This method is typically used in conjunction with Tag and MapTag to augment an existing set of
// tags or build a set of tags in an adhoc manner. For example:
//
//	var client aperture.Statsd
//	var condition bool
//
//	// Merge a set of adhoc tags with a base map of tags
//	base := map[string]interface{"foo": "bar"}
//	if condition {
//		client.Incr("name", lib.CombineTags(lib.MapTag(base), lib.Tag("baz", "qux")))
//	}
//
//	// Build a set of tags inline
//	client.Incr("name", lib.CombineTags(lib.Tag("foo", "bar")))
func CombineTags(tags ...Field) map[string]interface{} {
	combined := make(map[string]interface{})

	for _, tag := range tags {
		for _, it := range tag.items {
			combined[it.key] = it.value
		}
	}

	return combined
}
