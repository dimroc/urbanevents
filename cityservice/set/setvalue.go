// Modified set to also return an accompanying value
package set

import (
	"fmt"
	"sync"
)

type SetValue struct {
	m map[string]interface{}
	sync.RWMutex
}

type SetValueEntry struct {
	Key   string
	Value interface{}
}

func NewSetValue() *SetValue {
	return &SetValue{
		m: make(map[string]interface{}),
	}
}

func main2() {
	// Initialize our SetValue
	s := NewSetValue()

	data := struct {
		KeepAlive bool
		Message   string
	}{
		true,
		"luvin it",
	}

	data2 := struct {
		Message string
	}{
		"hatin it",
	}

	// Add example items
	s.Add("item1", data)
	s.Add("item1", data2)
	s.Add("item9", data)
	fmt.Printf("%d items\n", s.Len())

	// Clear all items
	s.Clear()
	if s.IsEmpty() {
		fmt.Printf("0 items\n")
	}

	s.Add("item2", data)

	// Check for existence
	if s.Has("item2") {
		fmt.Println("item2 does exist")
	}

	fmt.Println("list of all items:", s.List())

	// Remove some of our items
	s.Remove("item2")
	fmt.Println("list of all items:", s.List())
}

// Add add
func (s *SetValue) Add(item string, v interface{}) {
	s.Lock()
	defer s.Unlock()
	s.m[item] = v
}

// Remove deletes the specified item from the map
func (s *SetValue) Remove(item string) {
	s.Lock()
	defer s.Unlock()
	delete(s.m, item)
}

// Has looks for the existence of an item
func (s *SetValue) Has(item string) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[item]
	return ok
}

// Len returns the number of items in a set.
func (s *SetValue) Len() int {
	return len(s.List())
}

// Clear removes all items from the set
func (s *SetValue) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = make(map[string]interface{})
}

// IsEmpty checks for emptiness
func (s *SetValue) IsEmpty() bool {
	if s.Len() == 0 {
		return true
	}
	return false
}

// SetValue returns a slice of all items
func (s *SetValue) List() []SetValueEntry {
	s.RLock()
	defer s.RUnlock()
	list := make([]SetValueEntry, 0)
	for item := range s.m {
		list = append(list, SetValueEntry{Key: item, Value: s.m[item]})
	}
	return list
}

func (s *SetValue) ListAndClear() []SetValueEntry {
	s.Lock()
	defer s.Unlock()
	list := make([]SetValueEntry, 0)
	for item := range s.m {
		list = append(list, SetValueEntry{Key: item, Value: s.m[item]})
	}
	s.m = make(map[string]interface{})
	return list
}
