package worker

import (
	"encoding/json"
	"runtime"
)

// JobHandler is a callback function for handling job lifecycle events.
type JobHandler func(Job) error

// JobFunc is a function that processes a job and returns the result.
type JobFunc func(Job) ([]byte, error)

// The definition of the callback function.
type jobFunc struct {
	f       JobFunc
	timeout uint32
}

// Map for added function.
type jobFuncs map[string]*jobFunc

type systemInfo struct {
	GOOS, GOARCH, GOROOT, Version string
	NumCPU, NumGoroutine          int
	NumCgoCall                    int64
}

func SysInfo(job Job) ([]byte, error) {
	return json.Marshal(&systemInfo{
		GOOS:         runtime.GOOS,
		GOARCH:       runtime.GOARCH,
		GOROOT:       runtime.GOROOT(),
		Version:      runtime.Version(),
		NumCPU:       runtime.NumCPU(),
		NumGoroutine: runtime.NumGoroutine(),
		NumCgoCall:   runtime.NumCgoCall(),
	})
}

var memState runtime.MemStats

func MemInfo(job Job) ([]byte, error) {
	runtime.ReadMemStats(&memState)
	return json.Marshal(&memState)
}
