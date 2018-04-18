package itypes

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/emicklei/proto"
	"github.com/paralin/protods/generate"
	"github.com/pkg/errors"
	// "github.com/serenize/snaker"
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
func (g *Generator) GenerateCode(pf *proto.Proto) ([]byte, error) {
	var outp bytes.Buffer

	var packageName string
	var messages []*proto.Message
	for _, element := range pf.Elements {
		switch ele := element.(type) {
		case *proto.Package:
			packageName = ele.Name
		case *proto.Message:
			messages = append(messages, ele)
		}
	}

	if packageName == "" {
		return nil, errors.New("package name not found in proto file")
	}

	outp.WriteString("package ")
	outp.WriteString(packageName)
	outp.WriteString("\n\n")

	sort.Slice(messages, func(i int, j int) bool {
		return strings.Compare(messages[i].Name, messages[j].Name) == -1
	})

	for _, message := range messages {
		interName := "I" + message.Name

		// IHello is the interface type for Hello.
		outp.WriteString("// ")
		outp.WriteString(interName)
		outp.WriteString(" is the interface type for ")
		outp.WriteString(message.Name)
		outp.WriteString(".\n")

		// Hello is a hello message.
		if message.Comment != nil {
			outp.WriteString("// ")
			outp.WriteString(strings.TrimSpace(message.Comment.Message()))
			outp.WriteString("\n")
		}

		// type IHello interface {
		outp.WriteString("type ")
		outp.WriteString(interName)
		outp.WriteString(" interface {\n")

		for _, melement := range message.Elements {
			outp.WriteString(fmt.Sprintf("// %#v\n", melement))
				switch mele := melement.(type) {
				case *proto.Field:
					eleCamel := snaker.SnakeToCamel(mele.Name)
					outp.WriteString("\tGet")
					outp.WriteString(eleCamel)
					outp.WriteString("() ")
				}
			*/
		}
		outp.WriteString("}\n")
	}

	return outp.Bytes(), nil
}

func init() {
	generate.RegisterGenerator(generatorName, &Generator{})
}
