package main

import (
	"fmt"
	"os/exec"
	"time"
)

func printStates(states []map[string]interface{}) {
	out := make([]byte, len(states))
	for i, state := range states {
		if state == nil {
			out[i] = ' '
			continue
		}

		buildStatusIface, ok := state["build"]
		if !ok {
			out[i] = ' '
			continue
		}

		buildStatus, ok := buildStatusIface.(bool)
		if !ok {
			out[i] = 'e'
			continue
		}

		if buildStatus {
			out[i] = '-'
		} else {
			out[i] = 'X'
		}
	}
	fmt.Printf("%s\n", string(out))
}

func build() bool {
	cmd := exec.Command("go", "build")
	err := cmd.Run()
	return (err == nil)
}

func main() {
	ticker := time.Tick(10 * time.Second)
	tl := NewTimeline(30)
	go func() {
		for {
			states := <-tl.Updated()
			printStates(states)
		}
	}()
	for {
		go func() {
			buildStatus := build()
			tl.Update(func(state map[string]interface{}) {
				state["build"] = buildStatus
			})
		}()
		<-ticker
		tl.Tick()
	}
}
