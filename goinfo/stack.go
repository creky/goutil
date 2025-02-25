package goinfo

import (
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/gookit/goutil/basefn"
)

// some commonly consts
var (
	DefStackLen = 10000
	MaxStackLen = 100000
)

// GetCallStacks stacks is a wrapper for runtime.
// If all is true, Stack that attempts to recover the data for all goroutines.
//
// from glog package
func GetCallStacks(all bool) []byte {
	// We don't know how big the traces are, so grow a few times if they don't fit.
	// Start large, though.
	n := DefStackLen
	if all {
		n = MaxStackLen
	}

	// 4<<10 // 4 KB should be enough
	var trace []byte
	for i := 0; i < 10; i++ {
		trace = make([]byte, n)
		bts := runtime.Stack(trace, all)
		if bts < len(trace) {
			return trace[:bts]
		}
		n *= 2
	}
	return trace
}

// GetCallerInfo get caller func name and with base filename and line.
//
// returns:
//
//	github.com/gookit/goutil/goinfo_test.someFunc2(),stack_test.go:26
func GetCallerInfo(skip int) string {
	skip++ // ignore current func
	cs := GetCallersInfo(skip, skip+1)
	return basefn.FirstOr(cs, "")
}

// SimpleCallersInfo returns an array of strings containing
// the func name, file and line number of each stack frame leading.
func SimpleCallersInfo(skip, num int) []string {
	skip++ // ignore current func
	return GetCallersInfo(skip, skip+num)
}

// GetCallersInfo returns an array of strings containing
// the func name, file and line number of each stack frame leading.
//
// NOTICE: max should > skip
func GetCallersInfo(skip, max int) []string {
	var (
		pc         uintptr
		ok         bool
		line       int
		file, name string
	)

	callers := make([]string, 0, max-skip)
	for i := skip; i < max; i++ {
		pc, file, line, ok = runtime.Caller(i)
		if !ok {
			// The breaks below failed to terminate the loop, and we ran off the
			// end of the call stack.
			break
		}

		// This is a huge edge case, but it will panic if this is the case
		if file == "<autogenerated>" {
			break
		}

		fc := runtime.FuncForPC(pc)
		if fc == nil {
			break
		}

		if strings.ContainsRune(file, '/') {
			name = fc.Name()
			file = filepath.Base(file)
			// eg: github.com/gookit/goutil/goinfo_test.someFunc2(),stack_test.go:26
			callers = append(callers, name+"(),"+file+":"+strconv.Itoa(line))
		}

		// Drop the package
		// segments := strings.Split(name, ".")
		// name = segments[len(segments)-1]
	}

	return callers
}
