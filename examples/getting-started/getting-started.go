package gettingstarted

// NOTE: This file contains some simple tests, and is not yet a complete example.

type IExample interface{}
type IHello interface {
	GetSubject() string
	SetSubject(string)
	// NewMapField() IStringExampleMap
	GetMapField() IStringExampleMap
	SetMapField(IStringExampleMap)
}

type IStringExampleMap interface {
	Get(key string) IExample
	Set(key string, val IExample)
	ForEach(cb func(key string, val IExample) bool) bool
}
