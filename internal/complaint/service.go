package complaint

import "fmt"

// Service is the public API of the application: create complaints, fetch them
// and process them. It enforces validation and the status lifecycle.
type Service struct {
	store     *store
	processor Processor
}

// NewService returns a service with an empty in-memory store that uses p to
// process complaints.
func NewService(p Processor) *Service {
	return &Service{store: newStore(), processor: p}
}

// Create validates the input, stores a new complaint with status pending and
// returns a copy of it. Invalid input is rejected with ErrInvalidInput and
// nothing is stored.
func (s *Service) Create(in Input) (Complaint, error) {
	in = in.normalized()
	if err := in.Validate(); err != nil {
		return Complaint{}, err
	}
	c := s.store.add(in)
	return c.clone(), nil
}

// Get returns a copy of the complaint with the given ID, or ErrNotFound.
func (s *Service) Get(id ID) (Complaint, error) {
	c, err := s.store.get(id)
	if err != nil {
		return Complaint{}, err
	}
	return c.clone(), nil
}

// List returns copies of all complaints ordered by ID.
func (s *Service) List() []Complaint {
	return s.store.all()
}

// Process runs the processor on the complaint and moves it through the
// lifecycle: pending -> processing -> completed (or failed).
//
// Processing a complaint that is not pending returns ErrInvalidTransition and
// leaves it untouched. If the processor fails, the complaint becomes failed,
// has no result, and the processor's error is returned.
func (s *Service) Process(id ID) (Complaint, error) {
	c, err := s.store.get(id)
	if err != nil {
		return Complaint{}, err
	}

	if err := c.transition(StatusProcessing); err != nil {
		return Complaint{}, err
	}

	result, err := s.processor.Process(c.Input)
	if err != nil {
		if terr := c.transition(StatusFailed); terr != nil {
			return Complaint{}, terr
		}
		return c.clone(), fmt.Errorf("process complaint %d: %w", c.ID, err)
	}

	c.Result = &result
	if err := c.transition(StatusCompleted); err != nil {
		return Complaint{}, err
	}
	return c.clone(), nil
}
