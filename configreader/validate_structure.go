package configreader

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/kaptinlin/jsonschema"
)

// mustValidateConfigStructure takes an input document in as a sequence of
// bytes, for the purpose of checking it against a schema.
func mustValidateConfigStructure(documentVersion DocumentVersion, inputDocBuf *[]byte) {
	// Load the validation schema

	versionSchemaName := strings.ReplaceAll(string(documentVersion), ".", "-")
	file, err := os.Open(fmt.Sprintf("configreader/schemas/%s.yaml", versionSchemaName))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading validation schema: %v\n", err.Error())
		os.Exit(1)
	}

	validationSchemaYamlBuffer, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(1)
	}

	validationSchemaJsonBuffer := convertYamlBytesToJsonBytes(&validationSchemaYamlBuffer)

	compiler := jsonschema.NewCompiler()
	validationSchema, err := compiler.Compile(*validationSchemaJsonBuffer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error compiling schema: %v\n", err.Error())
		os.Exit(1)
	}

	inputDocJsonBytes := convertYamlBytesToJsonBytes(inputDocBuf)

	result := validationSchema.ValidateJSON(*inputDocJsonBytes)

	if !result.IsValid() {
		fmt.Fprint(os.Stderr, "Your input configuration is invalid :(\nHere's what the validation engine said:\n")
		detailedErrors := result.GetDetailedErrors()
		for field, message := range detailedErrors {
			fmt.Fprintf(os.Stderr, "%v: %v\n", field, message)
		}
		os.Exit(1)
	}
}
