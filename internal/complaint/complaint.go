package complaint

import (
	"fmt"
	"slices"
)

type ID int

type Category string

const (
	CategoryOrg       Category = "org"
	CategoryStaff     Category = "staff"
	CategoryCafeteria Category = "cafeteria"
	CategoryDormitory Category = "dormitory"
	CategoryOther     Category = "other"
)

type Validity string

const (
	ValidityValid   Validity = "VALID"
	ValidityReview  Validity = "REVIEW"
	ValidityInvalid Validity = "INVALID"
)


type Priority string

const (
	PriorityHigh   Priority = "HIGH"
	PriorityMedium Priority = "MEDIUM"
	PriorityLow    Priority = "LOW"
)

type Input struct {
	Category     Category
	Organization string
	School       string
	TeacherName  string
	Subject      string
	Message      string
}

type Result struct {
	Validity    Validity
	Priority    Priority
	ValidHits   int
	InvalidHits int
}

type Complaint struct {
	ID     ID
	Input  Input
	Status Status
	Result *Result
}

func (c *Complaint) transition(next Status) error {
	if !c.Status.CanTransitionTo(next) {
		return fmt.Errorf("%w: complaint %d cannot go from %s to %s",
			ErrInvalidTransition, c.ID, c.Status, next)
	}
	c.Status = next
	return nil
}

func (c *Complaint) clone() Complaint {
	cp := *c
	if c.Result != nil {
		r := *c.Result
		cp.Result = &r
	}
	return cp
}

func sortByID(list []Complaint) {
	slices.SortFunc(list, func(a, b Complaint) int { return int(a.ID) - int(b.ID) })
}
