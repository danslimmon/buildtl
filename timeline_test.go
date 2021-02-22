package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimeline_Basic(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	tl := NewTimeline(5)
	tl.Update(func(state map[string]interface{}) {
		state["hello"] = "world"
	})

	states := <-tl.Updated()
	assert.Equal(5, len(states), "timeline states array should not change length")
	valIface, ok := states[4]["hello"]
	assert.Equal(true, ok, "current state should have 'hello' key")
	val, ok := valIface.(string)
	assert.Equal(true, ok, "current state's 'hello' key should have a string value")
	assert.Equal("world", val, "current state's value for 'hello' key should be as specified")

	tl.Tick()
	states = <-tl.Updated()
	assert.Equal(5, len(states), "timeline states array should not change length")

	valIface, ok = states[4]["hello"]
	assert.Equal(true, ok, "current state should have 'hello' key")
	val, ok = valIface.(string)
	assert.Equal(true, ok, "current state's 'hello' key should have a string value")
	assert.Equal("world", val, "current state's value for 'hello' key should be as specified")

	valIface, ok = states[3]["hello"]
	assert.Equal(true, ok, "previous state should have 'hello' key")
	val, ok = valIface.(string)
	assert.Equal(true, ok, "previous state's 'hello' key should have a string value")
	assert.Equal("world", val, "previous state's value for 'hello' key should be as specified")
}
