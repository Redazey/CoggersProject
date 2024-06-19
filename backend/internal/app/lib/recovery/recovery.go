package recovery

import (
	"fmt"
	"runtime"
)

func Recovery() {
	if r := recover(); r != nil {
		pc, file, line, ok := runtime.Caller(2)
		if !ok {
			fmt.Println("Failed to recover")
			return
		}
		funcName := runtime.FuncForPC(pc).Name()
		fmt.Printf("Panic occurred at %s:%d (%s)n", file, line, funcName)
	}
}
