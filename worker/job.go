package worker

// Job represents a job that a worker can execute.
// It provides methods to interact with the job and send updates.
type Job interface {
	Err() error                                    // Returns any error associated with the job
	Data() []byte                                  // Returns the job data
	Fn() string                                    // Returns the function name
	SendWarning(data []byte)                       // Sends a warning message
	SendData(data []byte)                          // Sends data back to the client
	UpdateStatus(numerator, denominator int)       // Updates job progress
	Handle() string                                // Returns the job handle
	UniqueId() string                              // Returns the unique job identifier
}
