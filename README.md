# buildtl

`buildtl` is a running timeline of the build and test status of your code. You run it in a
one-row-high terminal and it produces output like this:

    ---ffffxxxxxxxxxxxxxxxxxxxxxxxfff-------------------------------fffffffffff-----------

It updates every 10 seconds, adding a new character at the right hand side to indicate the current
state. `x` means `go build` is failing, `f` means `go test ./...` is failing, and `-` means
everything's okay.
