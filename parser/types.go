package parser

// File represents a proto file.
type File struct {
	PackageName string
	Messages    []Message
	Maps        []Map
}

// Message is a known message type.
type Message struct {
	// Name is the name of the message.
	Name string
	// InterName is the name of the interface type.
	InterName string
	// Comment is the comment on the message.
	Comment string
	// Fields are the fields on the message.
	Fields []Field
}

// Field is a field in a message.
type Field struct {
	// Name is the snake_case name of the field.
	Name string
	// CamelName is the CamelCase name of the field.
	CamelName string
	// Comment is any comment on the field.
	Comment string
	// Type is the Go type of the field.
	Type string
	// Map indicates the field is a map type.
	Map *Map
}

// Map is a map type.
type Map struct {
	// Key is the key type for the map.
	// Original case
	Key string
	// Value is the value type for the map.
	// Original case
	Value string
	// ValuePtr is the value type without the I prefix.
	ValuePtr string
	// TypeName is the computed type name for the map interface.
	TypeName string
}
