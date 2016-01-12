package set

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// Returns true if setA and setB are equal
func EQS(a *S, b *S) bool {
	return a.EQ(b)
}

func TestSet(t *testing.T) {
	assert := assert.New(t)

	setA := ParseS("A B")
	setB := ParseS("A B C")
	setC := ParseS("B C D")

	assert.True(EQS(NewS("A", "B"), setA), "Did not create a new set correctly")
	assert.True(EQS(ParseS("A B"), setA), "Did not parse \"A B\" correctly")
	assert.True(NewS().Empty(), "The empty set must be empty")
	assert.Equal(setA.Size(), 2, "Size should be 2")
	assert.True(setA.IsSubsetOf(setB), "Should be a subset")
	assert.False(setB.IsSubsetOf(setA), "Should not be a subset")
	assert.True(setA.EQ(setA), "Should be equal to itself")
	assert.False(setA.EQ(setB), "Should not be equal")
	assert.True(setA.Contains("A"), "The set should contain the element A")
	assert.True(setA.Contains("B"), "The set should contain the element B")
	assert.False(setA.Contains("C"), "The set should not contain the element C")
	assert.True(EQS(setB.And(setC), ParseS("B C")), "The intersection should equal {B, C}")
	assert.True(EQS(setB.Or(setC), ParseS("A B C D")), "The union should equal {A, B, C, D}")
	assert.True(EQS(setC.Sub(setB), ParseS("D")), "Did not subtract sets correctly")
	assert.True(EQS(setB.Sub(setC), ParseS("A")), "Did not subtract sets correctly")
	assert.Equal(setC.String(), "{B, C, D}", "String conversion failed")

	setD := ParseS("A")

	// Insert
	setD.Insert("B")
	assert.True(EQS(setD, ParseS("A B")), "Failed to insert element B")

	setD.Insert("B")
	assert.True(EQS(setD, ParseS("A B")), "Inserting an existing element should do nothing")

	// Delete
	setD.Delete("A")
	assert.True(EQS(setD, ParseS("B")), "Failed to delete element A")

	setD.Delete("X")
	assert.True(EQS(setD, ParseS("B")), "Deleting a non-existent element should do nothing")

	// Elements
	elements := setA.Elements()
	assert.True(EQS(NewS(elements...), setA), "The two sets should be equal")

	// Map
	mapping := *setA.Map()
	assert.True(len(mapping) == 2 && mapping["A"] && mapping["B"], "Failed to extract the internal map")

	// Copy
	setE := setB.Copy()
	assert.True(EQS(setE, ParseS("A B C")), "Failed to copy")
}
