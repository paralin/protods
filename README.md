# Protobuf Datastructures

> Use any backing datastructure or datastore for complex Protobuf structures.

## Introduction

Allows developers to represent Protobuf objects with different data-structures. protods:

 - Generates interface types for Protobufs including field setters.
 - Declares common interfaces for various types of compatible backing data-structures.
 - Generates struct types binding generated protobuf structures to the backing data-structures.

Examples include:

 - Use a ctrie to allow efficient concurrent snapshotting of protobuf objects.
 - Use a key/value store to lazy-load protobufs from storage.
 - Load data from a Redis key/value store on-demand.
 
Note: this project is in the early development phase.
 
## Getting Started

To start getting a feel for how `protods` structures work, let's set up a simple project (under examples/getting-started).

First, fetch some dependencies:

```bash
go get -v github.com/golang/protobuf/protoc-gen-go
go get -v github.com/square/goprotowrap/cmd/protowrap
go get -v github.com/paralin/protods/cmd/protods
```

Create a protobuf structure:

```proto
// Hello is a hello message.
message Hello {
  // Subject is the subject of the Hello message.
  string subject = 1;
}
```

Generate the go code:

```bash
protoc --go_out=. ./getting-started.proto
```

Generate the protods code:

```bash
# Generate setters for all getters, interfaces for all types.
protods generate itypes getting-started.proto
```
 
## Code Generation Walkthrough

The "protods" tool is used to generate code depending on the desired output.

```proto
message Example {
   string str_field = 1;
   double num_field = 2;
   Example ex_field = 3;
   map<string, Example> map_field = 4;
}
```

Most modes of operation require generating the proto interface type:

```go
type IStringExampleMap interface {
	Get(key string) IExample
    Set(key string, val IExample)
    ForEach(cb func(key string, val IExample) bool) bool
}

type IExample interface {
   GetStrField() string
   SetStrField(string)
   GetNumField() float64
   SetNumField(float64)
   GetExField() IExample
   SetExField(IExample)
   // NewExField builds a new IObject of the same type as the parent.
   // For example, the generated New for proto types will return a proto type.
   NewExField() IExample
   NewMapField() IStringExampleMap
   GetMapField() IStringExampleMap
   // SetMapField sometimes requires a specific map type.
   // The proto generated types will use the given value if it is a map[]
   // Otherwise, they will clear the underlying map and copy the values with ForEach.
   SetMapField(IStringExampleMap)
}
```

The default generated Proto types implement half of the equation, the getters (GetStrField).

To make the generated Go message types compatible with the generated interfaces, setters are necessary:

```go
func (m *Example) SetStrField(val string) {
	if m != nil {
		m.StrField = val
	}
}
```

Now, an alternative data-structure might also satisfy `IExample`:

```go
// KeyValueExample is interchangeable at runtime with the proto object.
// This is because both implement the IExample type.
type KeyValueExample struct {
	m map[string]interface{}
}

// NewKeyValueExample returns a new IExample with a map backing it.
func NewKeyValueExample() IExample {
	return &KeyValueExample{m: make(map[string]interface{})}
}

// SetStrField sets the string field on the object.
func (m *KeyValueExample) SetStrField(val string) {
	if m != nil {
		m.m["str_field"] = val
	}
}

// GetStrField gets the string field from the object.
func (m *KeyValueExample) GetStrField() string {
	if m == nil {
		return ""
    }

	val, ok := m.m["str_field"]
	if !ok {
		return ""
	}
    
	return val.(string)
}
```

These are just examples of the types of structures that can be generated by the `protods` tool.

## Types of Backing Stores

This section describes the implemented types of backing stores for proto objects.

### Key/Value Store

A key/value store satisfies the following interface:

```go
// KeyValue contains values for a protobuf in a K/V store.
// Keys are generated using the field name.
type KeyValue interface {
	// Set stores the value for the key.	
	Set(key string, value interface{})
	// Get returns the value for the key.
	Get(key string) (bool, interface{})
	// Delete removes the value for the key.
	Delete(key string)
}
```

Given a key-value backed object:

```
message Upper {
	string id = 1;
	Lower lower = 2;
	map<string, Lower> lower_kv = 3;
}

message Lower {
	string value = 1;
}
```

The generated Key/Value backed code will generate keys like:

 - `upper.GetLower().GetValue()` -> `store.Get("/lower/value")`
 - `upper.GetLowerKv().Get("test").GetValue()` -> `store.Get("/lower_kv/test/value")`
 - `upper.GetId()` -> `store.Get("/id")`
