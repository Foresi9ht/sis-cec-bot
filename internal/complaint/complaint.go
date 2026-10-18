package complaint

import (
	"fmt"
	"slices"
)

// ID identifies a complaint. IDs are assigned by the store, start at 1 and are
// never reused.
type ID int

// Category is the type of complaint chosen by the student.
type Category string

const (
	CategoryOrg       Category = "org"       // student organizations
	CategoryStaff     Category = "staff"     // staff and teachers
	CategoryCafeteria Category = "cafeteria" // outside CEC jurisdiction unless a rights violation is described
	CategoryDormitory Category = "dormitory" // same as above
	CategoryOther     Category = "other"
)

// Validity is the first part of the processing result: does the complaint
// belong to the CEC's jurisdiction?
type Validity string

const (
	ValidityValid   Validity = "VALID"   // relevant keywords found
	ValidityReview  Validity = "REVIEW"  // nothing conclusive, needs a human
	ValidityInvalid Validity = "INVALID" // operational matter (food, plumbing, ...), not a CEC case
)

// Priority is the second part of the processing result: how urgent is it?
type Priority string

const (
	PriorityHigh   Priority = "HIGH"
	PriorityMedium Priority = "MEDIUM"
	PriorityLow    Priority = "LOW"
)

// Input is the domain-specific payload of a complaint. Which fields are
// required depends on the category (see Validate).
type Input struct {
	Category     Category
	Organization string // required for CategoryOrg
	School       string // required for CategoryStaff
	TeacherName  string // required for CategoryStaff
	Subject      string // optional, CategoryStaff only
	Message      string // always required
}

// Result is what processing produces. ValidHits and InvalidHits are the numbers
// of distinct keywords that decided the validity, kept so a result can be explained.
type Result struct {
	Validity    Validity
	Priority    Priority
	ValidHits   int
	InvalidHits int
}

// Complaint is one task of the application: input, lifecycle status and result.
type Complaint struct {
	ID     ID
	Input  Input
	Status Status
	Result *Result // nil until the complaint is completed
}

// transition moves the complaint to a new status or returns ErrInvalidTransition.
// Pointer receiver: it changes the complaint.
func (c *Complaint) transition(next Status) error {
	if !c.Status.CanTransitionTo(next) {
		return fmt.Errorf("%w: complaint %d cannot go from %s to %s",
			ErrInvalidTransition, c.ID, c.Status, next)
	}
	c.Status = next
	return nil
}

// clone returns a copy that shares no memory with c. Result is a pointer, so
// a plain struct copy would still point to the same Result.
func (c *Complaint) clone() Complaint {
	cp := *c
	if c.Result != nil {
		r := *c.Result
		cp.Result = &r
	}
	return cp
}

// sortByID orders complaints by ascending ID.
func sortByID(list []Complaint) {
	slices.SortFunc(list, func(a, b Complaint) int { return int(a.ID) - int(b.ID) })
}
