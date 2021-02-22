package main

import (
	"sync"
)

type Timeline struct {
	states  []map[string]interface{}
	mu      sync.Mutex
	updated chan []map[string]interface{}
}

func (tl *Timeline) Update(fn func(state map[string]interface{})) {
	tl.mu.Lock()
	defer tl.mu.Unlock()
	fn(tl.states[len(tl.states)-1])
	tl.updated <- tl.states
}

func (tl *Timeline) Tick() {
	tl.mu.Lock()
	defer tl.mu.Unlock()

	curState := tl.states[len(tl.states)-1]
	newState := make(map[string]interface{})
	for k, v := range curState {
		newState[k] = v
	}
	tl.states = append(tl.states[1:], newState)
	tl.updated <- tl.states
}

func (tl *Timeline) Updated() chan []map[string]interface{} {
	return tl.updated
}

func NewTimeline(length int) *Timeline {
	states := make([]map[string]interface{}, length)
	for i := range states {
		states[i] = make(map[string]interface{})
	}

	return &Timeline{
		states:  states,
		updated: make(chan []map[string]interface{}, 1),
	}
}
