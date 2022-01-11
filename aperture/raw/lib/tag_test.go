package lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestField(t *testing.T) {
	t.Parallel()

	field := Field{items: []item{}}
	assert.Equal(t, "", field.String())

	field = Field{items: []item{{"foo", 3}, {"bar", "baz"}}}
	assert.Equal(t, "foo", field.items[0].key)
	assert.Equal(t, 3, field.items[0].value)
	assert.Equal(t, "bar", field.items[1].key)
	assert.Equal(t, "baz", field.items[1].value)
	assert.Equal(t, "[foo: 3] [bar: baz]", field.String())
}

func TestNewTag(t *testing.T) {
	t.Parallel()

	tag := Tag("foo", "bar")
	assert.Equal(t, 1, len(tag.items))
	assert.Equal(t, item{"foo", "bar"}, tag.items[0])
}

func TestNewMapTag(t *testing.T) {
	t.Parallel()

	tag := MapTag(nil)
	assert.Equal(t, 0, len(tag.items))

	tag = MapTag(map[string]interface{}{})
	assert.Equal(t, 0, len(tag.items))

	tag = MapTag(map[string]interface{}{"foo": 3, "bar": "baz"})
	assert.Equal(t, 2, len(tag.items))
	assert.Contains(t, tag.items, item{"foo", 3})
	assert.Contains(t, tag.items, item{"bar", "baz"})
}

func TestCombineTags(t *testing.T) {
	t.Parallel()

	tags := CombineTags()
	assert.Equal(t, 0, len(tags))

	tags = CombineTags(Tag("foo", "bar"))
	assert.Equal(t, map[string]interface{}{"foo": "bar"}, tags)

	tags = CombineTags(Tag("foo", "bar"), Tag("bar", "baz"))
	assert.Equal(t, map[string]interface{}{"foo": "bar", "bar": "baz"}, tags)

	tags = CombineTags(MapTag(map[string]interface{}{"foo": "bar"}))
	assert.Equal(t, map[string]interface{}{"foo": "bar"}, tags)

	tags = CombineTags(MapTag(map[string]interface{}{"foo": "bar"}), Tag("bar", "baz"))
	assert.Equal(t, map[string]interface{}{"foo": "bar", "bar": "baz"}, tags)

	// Value override behavior
	tags = CombineTags(MapTag(map[string]interface{}{"foo": "bar"}), Tag("foo", "baz"))
	assert.Equal(t, map[string]interface{}{"foo": "baz"}, tags)

	// Value override behavior
	tags = CombineTags(Tag("foo", "baz"), MapTag(map[string]interface{}{"foo": "bar"}))
	assert.Equal(t, map[string]interface{}{"foo": "bar"}, tags)
}
