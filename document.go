package gonetable

// Implement Document interface for the structs you want to store to the DDB table.
//
// There are two required methods:
//
//	Gonetable_Key() returns the key that uniquely identifies the document.
//	Gonetable_TypeID() returns a string that specifies the type of the document.
//
// You can specify additional indeces with methods that return composite
// keys for them. They must be named with following pattern
//
//	Gonetable_[Index]Key() returns composite key for Index
type Document interface {
	Gonetable_Key() CompositeKey
	Gonetable_TypeID() string
}
