package complaint

import "fmt"

// store keeps all complaints in memory. It is not safe for concurrent use:
// Checkpoint 1 is synchronous. Concurrency is added in a later checkpoint.
//
// The map holds pointers (*Complaint) so that Service can change the status
// and result of a stored complaint in place. Everything that leaves the store
// is a copy, so callers cannot modify the stored data by accident.
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

// add stores a new pending complaint and assigns it the next free ID.
func (s *store) add(in Input) *Complaint {
	c := &Complaint{ID: s.nextID, Input: in, Status: StatusPending}
	s.complaints[c.ID] = c
	s.nextID++
	return c
}

// get returns the stored complaint (not a copy) or ErrNotFound.
func (s *store) get(id ID) (*Complaint, error) {
	c, ok := s.complaints[id]
	if !ok {
		return nil, fmt.Errorf("%w: id %d", ErrNotFound, id)
	}
	return c, nil
}

// all returns copies of all complaints ordered by ID. Map iteration order is
// random in Go, so the slice is sorted to make the output stable.
func (s *store) all() []Complaint {
	list := make([]Complaint, 0, len(s.complaints))
	for _, c := range s.complaints {
		list = append(list, c.clone())
	}
	sortByID(list)
	return list
}
