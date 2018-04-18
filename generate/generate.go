package generate

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/emicklei/proto"
	"github.com/pkg/errors"
)

// Generator generates a particular code file from a proto file.
type Generator interface {
	// GenerateCode generates code given the input proto file.
	GenerateCode(*proto.Proto) ([]byte, error)
	// GetUsage returns a usage description of the generator.
	GetUsage() string
	// GetShortName returns a short name for the generator, used in filenames.
	GetShortName() string
}

var registeredGenerators = make(map[string]Generator)

// RegisterGenerator registers a generator.
func RegisterGenerator(name string, gen Generator) {
	registeredGenerators[name] = gen
}

// GetGenerator gets a generator that was previously registered.
func GetGenerator(name string) Generator {
	return registeredGenerators[name]
}

// ForEachGenerator iterates over the generators.
func ForEachGenerator(cb func(name string, gen Generator) bool) {
	for id, gen := range registeredGenerators {
		if !cb(id, gen) {
			return
		}
	}
}

// Generate uses files to generate the proto output.
func Generate(gen Generator, protoPath, outputPath string) error {
	// Open the protobuf
	protoFilename := path.Base(protoPath)
	if !strings.HasSuffix(protoFilename, ".proto") {
		return errors.Errorf("expected .proto suffix: %v", protoFilename)
	}

	protoBaseName := strings.TrimSuffix(protoFilename, ".proto")

	f, err := os.Open(protoPath)
	if err != nil {
		return err
	}
	defer f.Close()

	parser := proto.NewParser(f)
	parser.Filename(protoFilename)
	parsedProto, err := parser.Parse()
	if err != nil {
		return errors.Wrap(err, "parse proto")
	}
	_ = f.Close()

	generatedCode, err := gen.GenerateCode(parsedProto)
	if err != nil {
		return err
	}

	fmtSrc, err := format.Source(generatedCode)
	if err != nil {
		// return err
		fmtSrc = generatedCode
	}

	// write the output
	outputFile := path.Join(outputPath, fmt.Sprintf("%s.%s.go", protoBaseName, gen.GetShortName()))
	return ioutil.WriteFile(outputFile, fmtSrc, 0644)
}
