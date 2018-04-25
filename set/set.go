package set

import (
	"sort"
	"strings"
)

// S is a set of strings
type S struct {
	elements map[string]bool
}

// NewS creates a set of strings from a list of elements e.g. `fruits := set.NewS("apple", "pear", "banana")`
func NewS(elements ...string) *S {
	result := S{make(map[string]bool)}
	for _, e := range elements {
		result.elements[e] = true
	}
	return &result
}

// ParseS creates a set from a space-delimited string (e.g. "A B C")
func ParseS(s string) *S {
	tokens := strings.Fields(s)
	return NewS(tokens...)
}

// Insert inserts an element into a set
func (s *S) Insert(e string) {
	s.elements[e] = true
}

// Delete deletes an element from  set
func (s *S) Delete(e string) {
	delete(s.elements, e)
}

// Size returns the number of elements in the set
func (s *S) Size() int {
	return len(s.elements)
}

// Empty return true if the set is empty
func (s *S) Empty() bool {
	return (s.Size() == 0)
}

// Copy returns a new set that contains the same elements as the instance
func (s *S) Copy() *S {
	result := NewS()
	for k, v := range s.elements {
		result.elements[k] = v
	}
	return result
}

// IsSubsetOf returns true if the set instance is a (mathematical) subset of the supplied set
func (s *S) IsSubsetOf(rhs *S) bool {
	for s := range s.elements {
		_, present := rhs.elements[s]
		if !present {
			return false
		}
	}
	return true
}

// EQ returns true if the two sets are the same
func (s *S) EQ(rhs *S) bool {
	return s.IsSubsetOf(rhs) && rhs.IsSubsetOf(s)
}

// Contains return true if the set contains the element
func (s *S) Contains(e string) bool {
	_, present := s.elements[e]
	return present
}

// And returns the intersection of two sets
func (s *S) And(rhs *S) *S {
	result := NewS()
	for s := range s.elements {
		if rhs.Contains(s) {
			result.Insert(s)
		}
	}
	return result
}

// Intersection returns the intersection of two sets
func (s *S) Intersection(rhs *S) *S {
	return s.And(rhs)
}

// Or returns the union of two sets
func (s *S) Or(rhs *S) *S {
	result := NewS()
	for k := range s.elements {
		result.Insert(k)
	}
	for k := range rhs.elements {
		result.Insert(k)
	}
	return result
}

// Union returns the union of two sets
func (s *S) Union(rhs *S) *S {
	return s.Or(rhs)
}

// Sub returns the set of elements which are members of the first set but
// not the second set.
func (s *S) Sub(a *S) *S {
	result := NewS()
	for e := range s.elements {
		if !a.Contains(e) {
			result.Insert(e)
		}
	}
	return result
}

// String ...
func (s *S) String() string {

	// Sort for deterministic output
	elements := s.Elements()
	sortS(elements)

	result := "{"
	for i, e := range elements {
		if i > 0 {
			result += ", "
		}
		result += e
	}
	result += "}"
	return result
}

// Elements returns the elements of the set as an array
func (s *S) Elements() []string {
	result := make([]string, 0, len(s.elements))
	for e := range s.elements {
		result = append(result, e)
	}
	return result
}

// Map returns a pointer to the internal map of a set, just for
// that special hacker that lives inside every developer
func (s *S) Map() *map[string]bool {
	return &s.elements
}

/////////////
// Sorting //
/////////////

type byString []string

func (s byString) Len() int {
	return len(s)
}

func (s byString) Swap(i int, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byString) Less(i int, j int) bool {
	return s[i] < s[j]
}

func sortS(s []string) {
	sort.Sort(byString(s))
}
