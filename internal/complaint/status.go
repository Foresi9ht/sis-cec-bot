package complaint

// Status is the lifecycle state of a complaint.
type Status string

const (
	// StatusPending: the complaint was created and validated, but not processed yet.
	StatusPending Status = "pending"
	// StatusProcessing: the complaint is being classified right now.
	StatusProcessing Status = "processing"
	// StatusCompleted: processing finished successfully and the result is stored.
	StatusCompleted Status = "completed"
	// StatusFailed: the processor returned an error; there is no result.
	StatusFailed Status = "failed"
)

// allowedTransitions is the whole lifecycle in one place.
//
//	pending -> processing -> completed
//	                      \-> failed
//
// completed and failed are terminal: they have no outgoing transitions.
var allowedTransitions = map[Status][]Status{
	StatusPending:    {StatusProcessing},
	StatusProcessing: {StatusCompleted, StatusFailed},
}

// CanTransitionTo reports whether moving from s to next is allowed.
// Value receiver: Status is a small immutable string-based value.
func (s Status) CanTransitionTo(next Status) bool {
	for _, allowed := range allowedTransitions[s] {
		if allowed == next {
			return true
		}
	}
	return false
}

// IsTerminal reports whether no further transitions are possible.
func (s Status) IsTerminal() bool {
	return len(allowedTransitions[s]) == 0
}

func (s Status) String() string { return string(s) }
