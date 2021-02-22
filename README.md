# runningtl

`runningtl` is a "running timeline" with an ncurses-based visualization system.

You make one like this:

    tl := runningtl.New()

You update it like this:

    tl.Update(func(state map[string]interface{}) error {
      state["weather"] = getWeather()
      return nil
    })

At any given time, the timeline has a **current state**. When you call `Update`, it modifies the
current state. Then, when you call `tl.Tick()`, the current state gets replaced by a new current
state. If the number of states in the timeline is then greater than the timeline's length, the
oldest state is deleted.

Whenever the data in the timeline changes, the channel returned by `tl.Updated` will receive the new
state as a `[]map[string]interface{}`.
