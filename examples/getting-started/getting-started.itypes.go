package gettingstarted

// IStringExampleMap is the map type for map<string, IExample>
type IStringExampleMap interface {
	Get(key string) IExample
	Set(key string, val IExample)
	ForEach(cb func(key string, val IExample) bool) bool
}

// StringExampleMap satisfies IStringExampleMap.
type StringExampleMap map[string]*Example

// Get returns a value from the map.
func (m StringExampleMap) Get(key string) IExample {
	if m == nil {
		return nil
	}
	return m[key]
}

// Set sets a value in the map.
func (m StringExampleMap) Set(key string, value IExample) {
	m[key] = value.(*Example)
}

// ForEach iterates over the map.
func (m StringExampleMap) ForEach(cb func(key string, val IExample) bool) bool {
	for k, v := range m {
		if !cb(k, v) {
			return false
		}
	}

	return true
}

// IExample is the interface type for Example.
type IExample interface {
}

func (m *Example) ToIExample() IExample {
	return (IExample)(m)
}

// IHello is the interface type for Hello.
// Hello is a hello message.
type IHello interface {
	GetSubject() string
	SetSubject(val string)
	GetMapFieldInter() IStringExampleMap
	SetMapField(val IStringExampleMap)
	NewMapField() IStringExampleMap
}

func (m *Hello) ToIHello() IHello {
	return (IHello)(m)
}

func (m *Hello) SetSubject(val string) {
	m.Subject = val
}

// _ is a type assertion
var _ IHello = &Hello{}

func (m *Hello) NewMapField() IStringExampleMap {
	return &StringExampleMap{}
}

func (m *Hello) GetMapFieldInter() IStringExampleMap {
	return StringExampleMap(m.GetMapField())
}

func (m *Hello) SetMapField(val IStringExampleMap) {
	m.MapField = (map[string]*Example)(val.(StringExampleMap))
}

// _ is a type assertion
var _ IHello = &Hello{}
