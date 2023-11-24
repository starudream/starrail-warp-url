package main

import (
	"runtime"
)

//goland:noinspection GoBoolExpressions
func isWindows() bool {
	return runtime.GOOS == "windows"
}
