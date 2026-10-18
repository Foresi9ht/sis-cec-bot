package complaint

import "fmt"

type Service struct {
	store     *store
	processor Processor
}

func NewService(p Processor) *Service {
	return &Service{store: newStore(), processor: p}
}

func (s *Service) Create(in Input) (Complaint, error) {
	in = in.normalized()
	if err := in.Validate(); err != nil {
		return Complaint{}, err
	}
	c := s.store.add(in)
	return c.clone(), nil
}

func (s *Service) Get(id ID) (Complaint, error) {
	c, err := s.store.get(id)
	if err != nil {
		return Complaint{}, err
	}
	return c.clone(), nil
}

func (s *Service) List() []Complaint {
	return s.store.all()
}

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
