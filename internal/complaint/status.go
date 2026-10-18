package complaint

// Status is the lifecycle state of a complaint.
type Status string

const (
	StatusPending Status = "pending"
	StatusProcessing Status = "processing"
	StatusCompleted Status = "completed"
	StatusFailed Status = "failed"
)

var allowedTransitions = map[Status][]Status{
	StatusPending:    {StatusProcessing},
	StatusProcessing: {StatusCompleted, StatusFailed},
}

func (s Status) CanTransitionTo(next Status) bool {
	for _, allowed := range allowedTransitions[s] {
		if allowed == next {
			return true
		}
	}
	return false
}

func (s Status) IsTerminal() bool {
	return len(allowedTransitions[s]) == 0
}

func (s Status) String() string { return string(s) }
