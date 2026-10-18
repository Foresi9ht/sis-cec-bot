package complaint

import (
	"errors"
	"testing"
)

func TestStatus_CanTransitionTo(t *testing.T) {
	tests := []struct {
		name string
		from Status
		to   Status
		want bool
	}{
		{"pending to processing", StatusPending, StatusProcessing, true},
		{"processing to completed", StatusProcessing, StatusCompleted, true},
		{"processing to failed", StatusProcessing, StatusFailed, true},
		{"pending to completed skips processing", StatusPending, StatusCompleted, false},
		{"pending to failed", StatusPending, StatusFailed, false},
		{"completed to processing (repeat)", StatusCompleted, StatusProcessing, false},
		{"completed to pending", StatusCompleted, StatusPending, false},
		{"failed to processing", StatusFailed, StatusProcessing, false},
		{"processing to pending", StatusProcessing, StatusPending, false},
		{"unknown status", Status("bogus"), StatusProcessing, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.from.CanTransitionTo(tt.to); got != tt.want {
				t.Errorf("%s.CanTransitionTo(%s) = %v, want %v", tt.from, tt.to, got, tt.want)
			}
		})
	}
}

func TestStatus_IsTerminal(t *testing.T) {
	tests := []struct {
		status Status
		want   bool
	}{
		{StatusPending, false},
		{StatusProcessing, false},
		{StatusCompleted, true},
		{StatusFailed, true},
	}

	for _, tt := range tests {
		if got := tt.status.IsTerminal(); got != tt.want {
			t.Errorf("%s.IsTerminal() = %v, want %v", tt.status, got, tt.want)
		}
	}
}

func TestComplaint_transition(t *testing.T) {
	c := &Complaint{ID: 7, Status: StatusPending}

	if err := c.transition(StatusProcessing); err != nil {
		t.Fatalf("pending -> processing: unexpected error %v", err)
	}
	if c.Status != StatusProcessing {
		t.Fatalf("status = %s, want %s", c.Status, StatusProcessing)
	}

	err := c.transition(StatusPending)
	if !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("processing -> pending: error = %v, want ErrInvalidTransition", err)
	}
	if c.Status != StatusProcessing {
		t.Errorf("status changed after a rejected transition: %s", c.Status)
	}
}
