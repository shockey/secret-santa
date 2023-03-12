package config

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/qri-io/jsonschema"
)

// mustValidateConfigStructure takes an input document in as a sequence of
// bytes, for the purpose of checking it against a schema. This intentionally
// happens before we parse the document as YAML into our structures, so that
// nonsensical inputs can be caught early.
func mustValidateConfigStructure(inputDocBuf *[]byte) {
	// Load the validation schema
	file, err := os.Open("config/validator-schema.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(1)
	}

	validationSchemaYamlBuffer, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(1)
	}

	validationSchemaJsonBuffer := convertYamlBytesToJsonBytes(&validationSchemaYamlBuffer)

	validationSchema := &jsonschema.Schema{}

	err = json.Unmarshal(*validationSchemaJsonBuffer, &validationSchema)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(1)
	}

	// Convert the document we're validating, then validate
	inputDocJsonBytes := convertYamlBytesToJsonBytes(inputDocBuf)

	errs, err := validationSchema.ValidateBytes(context.Background(), *inputDocJsonBytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(1)
	}

	if len(errs) > 0 {
		fmt.Fprint(os.Stderr, "Your input configuration is invalid :(\nHere's what the validation engine said:\n")
		for i, ve := range errs {
			fmt.Fprintf(os.Stderr, "%v: %v\n", i, ve.Error())
		}
	}

}
