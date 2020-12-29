package service

// job command
type JobCommand struct {
	Mac     string `json:"mac"`
	TaskId  string `json:"taskId"`
	JobId   string `json:"jobId"`
	Command string `json:"command"`
}

// dispatcher struct
// stop: use dispatcher.Stop()
type Dispatcher struct {
	C chan struct{}
}
