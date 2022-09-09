package gonetable_test

import "github.com/juranki/gonetable"

// Invalid index has one GSI with too short name
type InvalidIndex struct {
	Name string
}

func (sd1 *InvalidIndex) Gonetable_TypeID() string { return "sd1" }
func (sd1 *InvalidIndex) Gonetable_Key() gonetable.CompositeKey {
	return gonetable.CompositeKey{
		HashSegments:  []string{"a", "b"},
		RangeSegments: []string{"a", "b"},
	}
}
func (sd1 *InvalidIndex) Gonetable_AKey() gonetable.CompositeKey {
	return gonetable.CompositeKey{
		HashSegments:  []string{"a", "b"},
		RangeSegments: []string{"a", "b"},
	}
}

// MinimalDoc is minimal valid document
type MinimalDoc struct {
	Name string
}

func (sd1 *MinimalDoc) Gonetable_TypeID() string { return "sd1" }
func (sd1 *MinimalDoc) Gonetable_Key() gonetable.CompositeKey {
	return gonetable.CompositeKey{
		HashSegments:  []string{"a", "b"},
		RangeSegments: []string{"a", "b"},
	}
}

// WithIndex is valid document with one GSI
type WithIndex struct {
	Name string
}

func (sd1 *WithIndex) Gonetable_TypeID() string { return "wi1" }
func (sd1 *WithIndex) Gonetable_Key() gonetable.CompositeKey {
	return gonetable.CompositeKey{
		HashSegments:  []string{"wi", sd1.Name},
		RangeSegments: []string{"wi"},
	}
}
func (sd1 *WithIndex) Gonetable_GSI1Key() gonetable.CompositeKey {
	return gonetable.CompositeKey{
		HashSegments:  []string{"wi", sd1.Name},
		RangeSegments: []string{"wi"},
	}
}
