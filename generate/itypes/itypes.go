package itypes

import (
	"bytes"

	"github.com/paralin/protods/generate"
	"github.com/paralin/protods/parser"
)

const generatorName = "itypes"

// Generator generates the interface types for protos.
type Generator struct{}

// GetUsage returns a usage description of the generator.
func (g *Generator) GetUsage() string {
	return "generates interface types and proto setters"
}

// GetShortName returns the short name of the generator.
func (g *Generator) GetShortName() string {
	return generatorName
}

// GenerateCode generates code given the input proto file.
func (g *Generator) GenerateCode(pf *parser.File) ([]byte, error) {
	var outp bytes.Buffer

	outp.WriteString("package ")
	outp.WriteString(pf.PackageName)
	outp.WriteString("\n")

	for _, mapt := range pf.Maps {
		typeName := mapt.TypeName

		// IKeyValueMap is the map type for map<key, value>.
		outp.WriteString("\n// ")
		outp.WriteString(typeName)
		outp.WriteString(" is the map type for map<")
		outp.WriteString(mapt.Key)
		outp.WriteString(", ")
		outp.WriteString(mapt.Value)
		outp.WriteString(">\n")

		// type IKeyValueMap interface {
		outp.WriteString("type ")
		outp.WriteString(typeName)
		outp.WriteString(" interface {\n")

		// Get()
		outp.WriteString("\tGet(key string) ")
		outp.WriteString(mapt.Value)
		outp.WriteString("\n")

		// Set()
		outp.WriteString("\tSet(key string, val ")
		outp.WriteString(mapt.Value)
		outp.WriteString(")\n")

		// ForEach()
		outp.WriteString("\tForEach(cb func(key string, val ")
		outp.WriteString(mapt.Value)
		outp.WriteString(") bool) bool\n")

		outp.WriteString("}\n")

		// KeyValueMap satisfies IKeyValueMap.
		typeNameSansi := typeName[1:]
		valueSansi := mapt.Value
		if valueSansi[0] == 'I' {
			valueSansi = valueSansi[1:]
			valueSansi = "*" + valueSansi
		}

		outp.WriteString("\n// ")
		outp.WriteString(typeNameSansi)
		outp.WriteString(" satisfies ")
		outp.WriteString(typeName)
		outp.WriteString(".\n")

		// type KeyValueMap map[key]value
		outp.WriteString("type ")
		outp.WriteString(typeNameSansi)
		outp.WriteString(" map[string]")
		outp.WriteString(valueSansi)
		outp.WriteString("\n")

		isPrim := mapt.Value[0] != 'I'

		// Get returns a value from the map.
		outp.WriteString("\n// Get returns a value from the map.\n")
		outp.WriteString("func (m ")
		outp.WriteString(typeNameSansi)
		outp.WriteString(") Get(key string) ")
		outp.WriteString(mapt.Value)
		outp.WriteString(" {\n")
		if !isPrim {
			outp.WriteString("\tif m == nil { return nil }\n")
		} else {
			outp.WriteString("\tif m == nil { return new(")
			outp.WriteString(mapt.Value)
			outp.WriteString(") }\n")
		}
		outp.WriteString("\treturn m[key]\n")
		outp.WriteString("}\n")

		// Set sets a value in the map.
		outp.WriteString("\n// Set sets a value in the map.\n")
		outp.WriteString("func (m ")
		outp.WriteString(typeNameSansi)
		outp.WriteString(") Set(key string, value ")
		outp.WriteString(mapt.Value)
		outp.WriteString(") {\n")
		if isPrim {
			outp.WriteString("\tm[key] = value")
		} else {
			outp.WriteString("\tm[key] = value.(*")
			outp.WriteString(mapt.Value[1:])
			outp.WriteString(")\n")
		}
		outp.WriteString("}\n")

		// ForEach iterates over the map
		outp.WriteString("\n// ForEach iterates over the map.\n")
		outp.WriteString("func (m ")
		outp.WriteString(typeNameSansi)
		outp.WriteString(") ForEach(cb func(key string, val ")
		outp.WriteString(mapt.Value)
		outp.WriteString(") bool) bool {")
		outp.WriteString(`
	for k, v := range m {
		if !cb(k, v) {
			return false
		}
	}

	return true`)
		outp.WriteString("}\n")
	}

	for _, message := range pf.Messages {
		interName := message.InterName

		// IHello is the interface type for Hello.
		outp.WriteString("\n// ")
		outp.WriteString(interName)
		outp.WriteString(" is the interface type for ")
		outp.WriteString(message.Name)
		outp.WriteString(".\n")

		// Hello is a hello message.
		if message.Comment != "" {
			outp.WriteString("// ")
			outp.WriteString(message.Comment)
			outp.WriteString("\n")
		}

		// type IHello interface {
		outp.WriteString("type ")
		outp.WriteString(interName)
		outp.WriteString(" interface {\n")

		for _, field := range message.Fields {
			// GeT()
			outp.WriteString("\tGet")
			outp.WriteString(field.CamelName)
			if field.Map != nil {
				outp.WriteString("Inter")
			}
			outp.WriteString("() ")

			var typeName string
			if field.Map != nil {
				typeName = field.Map.TypeName
			} else {
				typeName = field.Type
			}

			outp.WriteString(typeName)
			outp.WriteString("\n")

			// Set()
			outp.WriteString("\tSet")
			outp.WriteString(field.CamelName)
			outp.WriteString("(val ")
			outp.WriteString(typeName)
			outp.WriteString(")\n")

			// New()
			if typeName[0] == 'I' {
				outp.WriteString("\tNew")
				outp.WriteString(field.CamelName)
				outp.WriteString("() ")
				outp.WriteString(typeName)
				outp.WriteString("\n")
			}
		}
		outp.WriteString("}\n")

		// Furthermore, augment the auto-generated proto types.
		for _, field := range message.Fields {
			var typeName string
			if field.Map != nil {
				typeName = field.Map.TypeName
			} else {
				typeName = field.Type
			}

			// func (m *Hello) NewSubject() ISubject
			if typeName[0] == 'I' {
				outp.WriteString("\nfunc (m *")
				outp.WriteString(message.Name)
				outp.WriteString(") New")
				outp.WriteString(field.CamelName)
				outp.WriteString("() ")
				outp.WriteString(typeName)
				outp.WriteString(" {\n\treturn &")
				outp.WriteString(typeName[1:])
				outp.WriteString("{}\n}\n")
			}

			// map type is generated at the beginning.
			if field.Map != nil {
				// generate funcs to get/set/new the map type

				// func (m *Hello) GetMapFieldInter() IMapFieldInter {
				outp.WriteString("\nfunc (m *")
				outp.WriteString(message.Name)
				outp.WriteString(") Get")
				outp.WriteString(field.CamelName)
				outp.WriteString("Inter() ")
				outp.WriteString(typeName)
				outp.WriteString(" {\n")

				typeNameWithouti := typeName
				if typeNameWithouti[0] == 'I' {
					typeNameWithouti = typeNameWithouti[1:]
				}

				// return IMapFieldInter(m.GetMapField())
				outp.WriteString("\treturn ")
				outp.WriteString(typeNameWithouti)
				outp.WriteString("(m.Get")
				outp.WriteString(field.CamelName)
				outp.WriteString("())\n}\n")
			}

			// func (m *Hello) SetSubject(val string)
			outp.WriteString("\nfunc (m *")
			outp.WriteString(message.Name)
			outp.WriteString(") Set")
			outp.WriteString(field.CamelName)
			outp.WriteString("(val ")

			outp.WriteString(typeName)
			outp.WriteString(") {\n")

			outp.WriteString("\tm.")
			outp.WriteString(field.CamelName)
			outp.WriteString(" = ")
			if field.Map == nil {
				outp.WriteString("val")
			} else {
				outp.WriteString("(map[")
				outp.WriteString(field.Map.Key)
				outp.WriteString("]")
				outp.WriteString(field.Map.ValuePtr)
				outp.WriteString(")")
				outp.WriteString("(val.(")
				outp.WriteString(field.Map.TypeName[1:])
				outp.WriteString("))")
			}
			outp.WriteString("\n}\n")

			// _ is a type assertion
			outp.WriteString("\n// _ is a type assertion\n")
			outp.WriteString("var _ ")
			outp.WriteString(message.InterName)
			outp.WriteString(" = &")
			outp.WriteString(message.Name)
			outp.WriteString("{}\n")
		}
	}

	return outp.Bytes(), nil
}

func init() {
	generate.RegisterGenerator(generatorName, &Generator{})
}
