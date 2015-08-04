package ideamart

import (
	"container/list"
	"log"
	"sync"
)

// An in-memory seession store for use with the USSD client.
// Automatically manages "garbage collection" by keeping track of updates.
// maxSize is the maximum number of sessions it will store before discarding the last updated session.
type inMemorySessionStore struct {
	maxSize     int
	currentSize int
	gcList      list.List
	lock        sync.Mutex
	store       map[string]*list.Element
}

func (s *inMemorySessionStore) deleteOldestIfFull() {
	log.Print(s.currentSize)
	if s.currentSize >= s.maxSize {
		e := s.gcList.Back()
		session := e.Value.(*USSDSession)
		delete(s.store, session.ID)
		s.gcList.Remove(e)
		s.currentSize--
	}
}

// Gets an USSD session by id. Returns nil if none exists.
func (s *inMemorySessionStore) Get(id string) *USSDSession {
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

// Saves a new USSD session. Removes the "oldest" to make space if insufficient.
func (s *inMemorySessionStore) Save(session USSDSession) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.deleteOldestIfFull()
	if s.store[session.ID] == nil {
		s.currentSize++
	}
	s.gcList.PushFront(&session)
	s.store[session.ID] = s.gcList.Front()
}

func (s *inMemorySessionStore) MaxSize() int {
	return s.maxSize
}

// Returns a properly initialized in-memory sessions store.
// maxSize is the maximum number of sessions it will store before discarding the "oldest" session.
func NewInMemorySessionStore(maxSize int) inMemorySessionStore {
	return inMemorySessionStore{maxSize: maxSize, store: map[string]*list.Element{}}
}
