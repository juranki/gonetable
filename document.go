package gonetable

// Implement Document interface for the structs you want to store to the DDB table.
//
// There are two required methods:
//
//	Gonetable_Key() returns the key that uniquely identifies the document.
//	Gonetable_TypeID() returns a string that specifies the type of the document.
//
// TODO:
// Additionally you can specify other functions with Gonetable_ -prefix to
// specify additional indexes, and fields that are computed when
// document is saved:
//
//	Gonetable_[Indexname]Key() returns composite key for additional index
//	Gonetable_Computed[Fieldname]() returns value for a computed field
type Document interface {
	Gonetable_Key() CompositeKey
	Gonetable_TypeID() string
}
