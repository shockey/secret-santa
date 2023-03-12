package config

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/icza/dyno"
	"github.com/qri-io/jsonschema"
	"gopkg.in/yaml.v2"
)

func mustValidateConfig(inputDocBuf *[]byte) {
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
		os.Exit(1)
	}

}

func convertYamlBytesToJsonBytes(yamlBuf *[]byte) *[]byte {
	var tmp interface{}

	err := yaml.Unmarshal(*yamlBuf, &tmp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(1)
	}

	// Context/author: https://stackoverflow.com/a/40737676
	tempMapS := dyno.ConvertMapI2MapS(tmp)

	jsonBytes, err := json.Marshal(tempMapS)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(1)
	}

	return &jsonBytes
}
