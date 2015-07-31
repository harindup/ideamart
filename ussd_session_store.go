package ideamart

import (
	"container/list"
	"log"
	"sync"
)

type InMemorySessionStore struct {
	MaxSize     int
	currentSize int
	gcList      list.List
	lock        sync.Mutex
	store       map[string]*list.Element
}

func (s *InMemorySessionStore) deleteOldestIfFull() {
	log.Print(s.currentSize)
	if s.currentSize >= s.MaxSize {
		e := s.gcList.Back()
		session := e.Value.(*USSDSession)
		delete(s.store, session.ID)
		s.gcList.Remove(e)
		s.currentSize--
	}
}

func (s *InMemorySessionStore) Get(id string) *USSDSession {
	s.lock.Lock()
	defer s.lock.Unlock()
	e := s.store[id]
	if e != nil {
		s.gcList.MoveToFront(e)
		if session, ok := e.Value.(*USSDSession); ok {
			return session
		} else {
			panic("Data Store corrupted: non-string key value in garbage collector")
			return nil
		}
	}
	return nil
}

func (s *InMemorySessionStore) Save(session USSDSession) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.deleteOldestIfFull()
	if s.store[session.ID] == nil {
		s.currentSize++
	}
	s.gcList.PushFront(&session)
	s.store[session.ID] = s.gcList.Front()
}

func NewInMemorySessionStore(maxSize int) InMemorySessionStore {
	return InMemorySessionStore{MaxSize: maxSize, store: map[string]*list.Element{}}
}
