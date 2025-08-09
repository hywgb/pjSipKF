package registrar

import (
	"context"
	"sync"
	"time"
)

type inMemory struct {
	mu       sync.RWMutex
	records  map[string][]Binding
}

func NewInMemory() Service {
	return &inMemory{records: make(map[string][]Binding)}
}

func (s *inMemory) Register(_ context.Context, user, contact string, expiresSeconds int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	exp := time.Now().Add(time.Duration(expiresSeconds) * time.Second).Unix()
	list := s.records[user]
	// replace if contact exists
	replaced := false
	for i := range list {
		if list[i].Contact == contact {
			list[i].ExpiresUnix = exp
			replaced = true
			break
		}
	}
	if !replaced {
		list = append(list, Binding{Contact: contact, ExpiresUnix: exp})
	}
	s.records[user] = list
	return nil
}

func (s *inMemory) Deregister(_ context.Context, user, contact string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	list := s.records[user]
	res := make([]Binding, 0, len(list))
	for _, b := range list {
		if b.Contact != contact {
			res = append(res, b)
		}
	}
	if len(res) == 0 {
		delete(s.records, user)
	} else {
		s.records[user] = res
	}
	return nil
}

func (s *inMemory) Lookup(_ context.Context, user string) ([]Binding, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := s.records[user]
	now := time.Now().Unix()
	res := make([]Binding, 0, len(list))
	for _, b := range list {
		if b.ExpiresUnix > now {
			res = append(res, b)
		}
	}
	return res, nil
}