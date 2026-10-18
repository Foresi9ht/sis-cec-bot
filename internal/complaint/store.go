package complaint

import "fmt"

type store struct {
	complaints map[ID]*Complaint
	nextID     ID
}

func newStore() *store {
	return &store{
		complaints: make(map[ID]*Complaint),
		nextID:     1,
	}
}

func (s *store) add(in Input) *Complaint {
	c := &Complaint{ID: s.nextID, Input: in, Status: StatusPending}
	s.complaints[c.ID] = c
	s.nextID++
	return c
}

func (s *store) get(id ID) (*Complaint, error) {
	c, ok := s.complaints[id]
	if !ok {
		return nil, fmt.Errorf("%w: id %d", ErrNotFound, id)
	}
	return c, nil
}

func (s *store) all() []Complaint {
	list := make([]Complaint, 0, len(s.complaints))
	for _, c := range s.complaints {
		list = append(list, c.clone())
	}
	sortByID(list)
	return list
}
