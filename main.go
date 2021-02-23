package main

import (
	"fmt"
	"os/exec"
	"time"

	tm "github.com/buger/goterm"
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
			out[i] = 'x'
			continue
		}

		testStatusIface, ok := state["test"]
		if !ok {
			out[i] = ' '
			continue
		}

		testStatus, ok := testStatusIface.(bool)
		if !ok {
			out[i] = 'e'
			continue
		}

		if testStatus {
			out[i] = '-'
		} else {
			out[i] = 'f'
			continue
		}
	}
	fmt.Printf("\n%s", string(out))
}

func build() bool {
	cmd := exec.Command("go", "build")
	err := cmd.Run()
	return (err == nil)
}

func test() bool {
	cmd := exec.Command("go", "test", "./...")
	err := cmd.Run()
	return (err == nil)
}

func termWidth() int {
	return tm.Width()
}

func main() {
	ticker := time.Tick(10 * time.Second)
	tl := NewTimeline(termWidth())
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

			testStatus := test()
			tl.Update(func(state map[string]interface{}) {
				state["test"] = testStatus
			})
		}()
		<-ticker
		tl.Tick()
	}
}
