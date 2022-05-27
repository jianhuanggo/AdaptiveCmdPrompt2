package main

import (
	"runtime"
	"Con_Utils/cmd"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cmd.Execute()
}
