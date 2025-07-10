package client

// StatusHandler is a callback function for handling job status updates.
// Parameters: handle, known, running, numerator, denominator
type StatusHandler func(string, bool, bool, uint64, uint64)

// Status represents the current status of a job.
type Status struct {
	Handle                 string // Job handle
	Known, Running         bool   // Status flags
	Numerator, Denominator uint64 // Progress information
}
